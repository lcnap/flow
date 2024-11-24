package main

import (
	"context"
	"io"
	config "lcnap/flow/config"
	"lcnap/flow/plugin"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

var defaultProxy *Proxy = &Proxy{
	client: http.DefaultClient,
}

type Proxy struct {
	mu      sync.Mutex
	Servers []*http.Server

	client *http.Client

	errlogger    *slog.Logger
	accesslogger *slog.Logger

	Conf config.Conf
}

func (p *Proxy) Init(conf *config.Conf) error {

	p.initLog(conf)
	p.Conf = *conf

	p.errlogger.Debug("copy conf.", "new", &p.Conf, "old", conf)

	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range conf.Server {
		s := p.initServer(v)
		p.Servers = append(p.Servers, s)

	}
	return nil
}

func (p *Proxy) initLog(conf *config.Conf) {
	p.accesslogger = NewLogger(conf.Log.Access)
	p.errlogger = NewLogger(conf.Log.Error)
}

func (p *Proxy) start() {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range p.Servers {
		go p.runServer(v)
	}

}

func (p *Proxy) runServer(v *http.Server) {
	p.errlogger.Info("server is running.", "server", v.Addr)
	err := v.ListenAndServe()
	if err == http.ErrServerClosed {
		p.errlogger.Warn("server is closing", "server", v.Addr)
	} else {
		p.errlogger.Warn("server started failed.", "server", v.Addr)
	}
}

func (p *Proxy) initServer(srv config.Server) *http.Server {

	return &http.Server{
		Addr:         srv.Listen,
		Handler:      p.initHandler(srv.Route),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
}

func (p *Proxy) initHandler(route []config.Route) http.Handler {
	sm := http.NewServeMux()
	plugin.ResetGoRuntime()
	for _, v := range route {
		if v.Location != "" && v.Pass != "" {
			sm.HandleFunc(v.Location, func(w http.ResponseWriter, r *http.Request) {

				req, err := http.NewRequest(r.Method, v.Pass+r.RequestURI, r.Body)
				req.Header = r.Header.Clone()

				if err != nil {
					p.errlogger.Warn(err.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				resp, err := defaultProxy.client.Do(req)
				if err != nil {
					p.errlogger.Warn(err.Error())
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				p.accesslogger.Info("proxy", "ip", r.RemoteAddr, "req", r.RequestURI, "pass", req.URL, "resp", resp.Status)
				for i, v := range resp.Header {
					w.Header().Add(i, v[0])
				}
				w.WriteHeader(resp.StatusCode)
				io.Copy(w, resp.Body)
				defer resp.Body.Close()
			})
			continue
		}
		if v.Location != "" && v.Handler != "" {
			handlerid := v.Handler
			handler, err := plugin.LoadGo(handlerid)
			if err != nil {
				p.errlogger.Error("load go src failed.", "err", err)
				continue
			}

			sm.HandleFunc(v.Location, func(w http.ResponseWriter, r *http.Request) {
				p.accesslogger.Info("execute plugin.", "ip", r.RemoteAddr, "req", r.RequestURI, "handler", handlerid)
				handler(w, r)
			})
		}
	}
	return sm
}

func (p *Proxy) stop() {
	p.end(func(s *http.Server) error {
		return s.Shutdown(context.Background())
	})

}

func (p *Proxy) quit() {
	p.end(func(s *http.Server) error {
		return s.Close()
	})
}

func (p *Proxy) end(f func(*http.Server) error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, v := range p.Servers {
		err := f(v)
		if err != nil {
			p.errlogger.Error(err.Error())
		}
	}
}

func (p *Proxy) reload(conf *config.Conf) {
	p.mu.Lock()
	defer p.mu.Unlock()

	newServers := []*http.Server{}
	m := map[*http.Server]bool{}
	for _, cv := range conf.Server {
		find := false
		for _, pv := range p.Servers {
			if cv.Listen == pv.Addr {
				pv.Handler = p.initHandler(cv.Route)
				m[pv] = true
				find = true
				newServers = append(newServers, pv)
			}
		}
		if !find {
			s := p.initServer(cv)
			go p.runServer(s)
			newServers = append(newServers, s)
		}

	}

	for _, v := range p.Servers {
		if !m[v] {
			err := v.Shutdown(context.Background())
			if err != nil {
				p.errlogger.Error(err.Error())
			}
		}
	}

	p.Servers = newServers
}

func (p *Proxy) reopen(conf *config.Conf) {
	p.initLog(conf)
}
