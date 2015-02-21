package weblog

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}

type WebLog struct {
}

func Handle(w http.ResponseWriter, r *http.Request) {

}
