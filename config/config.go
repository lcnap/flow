package config

type Conf struct {
	Server   []Server
	Upstream []Upstream
	Log      Log
}

type Server struct {
	Scheme string
	Listen string
	Ssl    interface{}
	Route  []Route
}

type Route struct {
	Location string
	Handler  string
	Pass     string
}

type Upstream struct {
	Name     string
	EndPoint string
}

type Log struct {
	Access string
	Error  string
}
