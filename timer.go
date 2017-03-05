package main

import (
	"fmt"
	"time"

	"golang.org/x/net/websocket"
)

type WorkLog struct {
	User      string
	BeginTime uint64
	EndTime   uint64
	Task      string
}

func doTimer(ws *websocket.Conn, m Message, cancel chan Message) {
	m.Text = "got it"
	postMessage(ws, m)
	go func(m Message) {
		fmt.Println(m.User)
		select {
		case _ = <-cancel:
			m.Text = "Your timer is canceled"
			postMessage(ws, m)
		case _ = <-time.After(30 * time.Second):
			m.Text = "Your time is up"
			postMessage(ws, m)
		}
	}(m)
}
