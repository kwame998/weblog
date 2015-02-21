package weblog

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func NewWebLog() *WebLog {
	return &WebLog{register: make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool)}
}

type WebLog struct {
	register    chan *connection
	unregister  chan *connection
	connections map[*connection]bool
}

func (wl *WebLog) Run() {
	for {
		select {
		case c := <-wl.unregister:
			wl.connections[c] = false
			delete(wl.connections, c)
			c.ws.Close()
		case c := <-wl.register:
			wl.connections[c] = true
			fmt.Println("Just registered a client")
		}
	}
}

func (wl *WebLog) Write(p []byte) (int, error) {
	for conn := range wl.connections {
		conn.send <- p
	}

	return len(p), nil
}

func (wl *WebLog) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	wl.register <- c
	defer func() { wl.unregister <- c }()

	c.send <- []byte("Welcome to weblog!")
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
