package main

import (
	"net/http"

	"github.com/paked/weblog"
)

var l *weblog.WebLogger = weblog.NewWebLogger()

func main() {
	http.HandleFunc("/ws", l.Handle)
	http.HandleFunc("/", home)

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	l.Println("Someone accessed your homepage...")
	http.ServeFile(w, r, "index.html")
}
