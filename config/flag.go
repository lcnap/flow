package config

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"os"
)

/*
proxy -c config.json
proxy -s reload,stop,quit,reopen
proxy -t
proxy -h
proxy -v
*/

func ParseFlag() FlagArg {
	var s SignalArg
	flag.Var(&s, "s", "send a signal to process. signals:reload,stop,quit,reopen")

	var c = flag.String("c", "./config.json", "configuration file.")
	var t = flag.Bool("t", false, "test the configuration file.")
	var e = flag.String("e", "./logs/error.log", "error file.")
	flag.Parse()

	return FlagArg{
		Singal:         s.String(),
		ConfigFile:     *c,
		TestConfigFile: *t,
		ErrorFile:      *e,
	}
}

type FlagArg struct {
	Singal         string
	ConfigFile     string
	TestConfigFile bool
	ErrorFile      string
}

func (a *FlagArg) LoadConfig(path string) (*Conf, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var config Conf
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type SignalArg string

func (s *SignalArg) String() string {
	return string(*s)
}

func (s *SignalArg) Set(value string) error {
	if value != "reload" && value != "stop" && value != "quit" && value != "reopen" {
		return errors.New("reload,stop,quit,reopen")
	}
	*s = SignalArg(value)
	return nil
}
