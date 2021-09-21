package server

import (
	"go-travel/chatroom/logic"
	"net/http"
)

func RegisterHandle() {
	go logic.Broadcaster.Start()
	http.HandleFunc("/", homeHandleFunc)
	http.HandleFunc("/ws", WebSocketHandleFunc)
}
