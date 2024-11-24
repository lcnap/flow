package main

import (
	"fmt"
	"io"
	config "lcnap/flow/config"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var mainlogger *slog.Logger = NewLogger("./logs/error.log")

var (
	USR1    = syscall.Signal(0xa)
	signals = map[string]syscall.Signal{
		"stop":   syscall.SIGQUIT,
		"quit":   syscall.SIGINT,
		"reload": syscall.SIGHUP,
		"reopen": USR1,
	}
)

func main() {
	arg := config.ParseFlag()

	if hasStarted() {
		//向 main 进程 发送 信号
		pid, err := getPid()
		if err != nil {
			mainlogger.Error("open pid file failed.", "err", err)
			return
		}
		proc, err := os.FindProcess(pid)
		if err != nil {
			mainlogger.Error("failed to find proccess.", "err", err)
			return
		}
		if sig, ok := signals[arg.Singal]; ok {
			err := proc.Signal(sig)
			if err != nil {
				mainlogger.Error("failed to send signal.", "err", err)
			}
		} else {
			mainlogger.Error("unsupported signal.", "signal", arg.Singal)
		}

		return
	}

	conf, err := arg.LoadConfig(arg.ConfigFile)
	if err != nil {
		mainlogger.Error(err.Error())
		return
	}
	defaultProxy.Init(conf)
	defaultProxy.start()

	savePid(os.Getpid())

	msigfunc := map[os.Signal]func(){
		syscall.SIGQUIT: func() {
			defaultProxy.stop()
			clearPid()
			os.Exit(0)
		},
		syscall.SIGINT: func() {
			defaultProxy.quit()
			clearPid()
			os.Exit(0)
		},
		syscall.SIGHUP: func() {
			conf, err = arg.LoadConfig(arg.ConfigFile)
			if err != nil {
				mainlogger.Error(err.Error())
				return
			}
			defaultProxy.reload(conf)
		},
		USR1: func() {
			conf, err = arg.LoadConfig(arg.ConfigFile)
			if err != nil {
				mainlogger.Error(err.Error())
				return
			}
			defaultProxy.reopen(conf)
		},
		/* syscall.SIGKILL: func() {
			clearPid()
			os.Exit(0)
		}, */
	}

	waitSignal(msigfunc)

}

const (
	PidPath string = "./logs/lcnap.pid"
)

func savePid(pid int) {
	f, err := os.OpenFile(PidPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		mainlogger.Error("write pid failed.", "err", err)
		return
	}
	defer f.Close()
	f.Seek(0, 0)
	_, err = io.WriteString(f, fmt.Sprintf("%d", pid))
	if err != nil {
		mainlogger.Error("write pid failed.", "err", err)
	}
}

func clearPid() {
	os.Remove(PidPath)
}

func getPid() (int, error) {
	f, err := os.Open(PidPath)
	if err != nil {
		return -1, err
	}
	defer f.Close()

	buf, err := io.ReadAll(f)
	if err != nil {
		return -1, err
	}
	spid := string(buf)
	pid, err := strconv.Atoi(spid)
	if err != nil {
		return -1, err
	}
	return pid, nil
}

func hasStarted() bool {
	pid, err := getPid()
	if err != nil {
		return false
	}
	if pid > 0 {
		return true
	}
	return false
}

func waitSignal(m map[os.Signal]func()) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan)
	for {
		//处理信号.
		s := <-sigchan
		if f, ok := m[s]; ok {
			mainlogger.Info("received signal.", "signal", s)
			f()
		}
	}

}
