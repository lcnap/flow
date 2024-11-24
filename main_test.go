package main

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
)

func Test_hasStarted(t *testing.T) {
	os.Remove(PidPath)
	f := hasStarted()
	if f {
		t.Fatal("except false.")
	}
	savePid(os.Getpid())
	f = hasStarted()
	if !f {
		t.Fatal("except true.")
	}
}

func Test_reload(t *testing.T) {
	var wg sync.WaitGroup

	wg.Add(1)
	s := http.Server{
		Addr: ":8080",
	}
	go func() {
		s.ListenAndServe()
		wg.Done()
	}()
	for i := 0; i < 1; i++ {

		go func(i int) {
			mux := http.ServeMux{}
			mux.HandleFunc(fmt.Sprintf("/bar%d", i), func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "Hello, %d", i)
			})
			s.Handler = &mux

		}(i)
	}
	wg.Wait()
}

func Test_waitSignal(t *testing.T) {
	msigfunc := map[os.Signal]func(){
		syscall.SIGQUIT: func() {
			t.Log(syscall.SIGQUIT)
		},
		syscall.SIGINT: func() {
			t.Log(syscall.SIGINT)
		},
		syscall.SIGHUP: func() {
			t.Log(syscall.SIGHUP)
		},
		USR1: func() {
			t.Log(USR1, "xx")
		},
		syscall.SIGKILL: func() {
			t.Log(syscall.SIGKILL)
		},
	}
	t.Log(os.Getpid())
	waitSignal(msigfunc)
}
