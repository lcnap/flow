package plugin

import (
	"net/http"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

type Handler func(http.ResponseWriter, *http.Request)

var _goruntime *interp.Interpreter = interp.New(interp.Options{})
var LoadedSrcMap map[string]bool = map[string]bool{}

func goruntime() *interp.Interpreter {
	return _goruntime
}

func LoadGo(path string) (Handler, error) {
	paths := strings.Split(path, ":")
	src, fun := paths[0], paths[1]
	if err := LoadGoSrc(src); err != nil {
		return nil, err
	}
	f, err := LoadGoFunc(fun)
	if err != nil {
		return nil, err
	}

	return Handler(f), nil

}

func LoadGoSrc(path string) error {
	if LoadedSrcMap[path] {
		return nil
	}
	goruntime().Use(stdlib.Symbols)
	_, err := goruntime().EvalPath(path)
	if err != nil {
		return nil
	}
	LoadedSrcMap[path] = true
	return nil
}

func LoadGoFunc(fun string) (Handler, error) {
	f, err := goruntime().Eval(fun)
	if err != nil {
		return nil, err
	}
	ff, ok := f.Interface().(func(http.ResponseWriter, *http.Request))
	if ok {
		return ff, nil
	}
	return nil, err
}

func ResetGoRuntime() {
	_goruntime = interp.New(interp.Options{})
	LoadedSrcMap = map[string]bool{}
}
