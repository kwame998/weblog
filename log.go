package weblog

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

func NewWebLog() WebLog {
	return WebLog{make(chan []byte, 256)}
}

type WebLog struct {
	send chan []byte
}

func (wl WebLog) Write(p []byte) (int, error) {
	wl.send <- p
	return len(p), nil
}

func (wl WebLog) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	ws.WriteMessage(websocket.TextMessage, []byte("Welcome to weblog!"))

	go func() {
		for message := range wl.send {
			err := ws.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error: ", err)
				break
			}
			fmt.Println("SENDING MESSAGE")
		}

		ws.Close()
	}()
}
