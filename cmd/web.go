package main

import (
	"net/http"

	"github.com/paked/weblog"
	"log"
)

var l *log.Logger

func main() {
	wl := weblog.NewWebLog()
	http.HandleFunc("/ws", wl.Handle)
	http.HandleFunc("/", home)

	l = log.New(wl, "LOGGING: ", log.Ldate|log.Ltime)
	go wl.Run()
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	l.Println("Someone accessed your homepage...")
	http.ServeFile(w, r, "index.html")
}
