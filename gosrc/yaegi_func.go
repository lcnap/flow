package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func fw(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RequestURI + ";msg from gosrc:f func"))
}
func f(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(r.RequestURI + ";msg from gosrc:fww func"))
}
func c(msg string) string {
	fmt.Println(msg)
	return time.Now().String()
}

type Msg struct {
	Content string
	Id      string
	Time    time.Time
}

func longfunc(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Msg{
		Content: "content",
		Id:      "id",
		Time:    time.Now(),
	})
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}

type Pong struct {
	Ping string
}

func ping(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(Pong{
		"Pong",
	})
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
