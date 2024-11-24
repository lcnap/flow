package main

import (
	"fmt"
	"testing"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func Test_runYaegi(t *testing.T) {
	i := interp.New(interp.Options{})

	i.Use(stdlib.Symbols)
	_, err := i.EvalPath("./yaegi_func.go")
	if err != nil {
		panic(err)
	}

	f, err := i.Eval(`main.c`)
	if err != nil {
		panic(err)
	}
	rf := f.Interface().(func(string) string)
	r := rf("msg from go")
	fmt.Println(r)
}

func Test_runYaegi1(t *testing.T) {
	i := interp.New(interp.Options{})
	i.Use(stdlib.Symbols)

	pro, err := i.CompilePath("./yaegi_func.go")
	if err != nil {
		panic(err)
	}

	_, err = i.Execute(pro)
	if err != nil {
		panic(err)
	}

}
