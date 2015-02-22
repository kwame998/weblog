package weblog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func NewWebLogger() *WebLogger {
	ww := &webWriter{register: make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool)}

	l := log.New(ww, "", log.Lshortfile)
	go ww.Run()

	return &WebLogger{l, ww}
}

type WebLogger struct {
	*log.Logger

	writer *webWriter
}

func (wl *WebLogger) Handle(w http.ResponseWriter, r *http.Request) {
	wl.writer.Handle(w, r)
}

type webWriter struct {
	register    chan *connection
	unregister  chan *connection
	connections map[*connection]bool
}

func (ww *webWriter) Run() {
	for {
		select {
		case c := <-ww.unregister:
			ww.connections[c] = false
			delete(ww.connections, c)
			c.ws.Close()
		case c := <-ww.register:
			ww.connections[c] = true
			fmt.Println("Just registered a client")
		}
	}
}

func (ww *webWriter) Write(p []byte) (int, error) {
	for conn := range ww.connections {
		conn.send <- p
	}

	return len(p), nil
}

func (ww *webWriter) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	ww.register <- c
	defer func() { ww.unregister <- c }()

	c.send <- []byte("Welcome to the weblog!")
	c.writer()
}

type connection struct {
	send chan []byte
	ws   *websocket.Conn
}

func (c connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
