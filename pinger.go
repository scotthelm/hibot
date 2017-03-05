package main

import (
	"sync/atomic"
	"time"

	"golang.org/x/net/websocket"
)

func ping(ws *websocket.Conn) {
	select {
	case _ = <-time.After(5 * time.Second):
		id := atomic.AddUint64(&counter, 1)
		m := Message{Id: id, Type: "ping"}
		postMessage(ws, m)
		ping(ws)
	}
}
