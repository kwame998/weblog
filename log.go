package weblog

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type WebLog struct {
	send chan []byte
}

func (wl *WebLog) Write(p []byte) (int, error) {
	wl.send <- p
	return len(p), nil
}

func (wl WebLog) Handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	for message := range wl.send {
		err := ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}

	ws.Close()
}
