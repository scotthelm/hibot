package main

import (
	"fmt"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type WorkLog struct {
	User      string
	BeginTime int64
	EndTime   int64
	Duration  uint64
	Task      string
}

func doTimer(ws *websocket.Conn, m Message, cancel chan Message) {
	mutex.Lock()
	worklog, found := state[m.User]
	mutex.Unlock()
	if found {
		m.Text = fmt.Sprintf("You are already working on %s", worklog.Task)
		postMessage(ws, m)
	} else {
		addState(m)
		go func(m Message) {
			fmt.Println(m.User)
			select {
			case _ = <-cancel:
				m.Text = "Your timer is canceled"
				postMessage(ws, m)
				removeState(m)
			case _ = <-time.After(30 * time.Second):
				m.Text = "Your time is up"
				postMessage(ws, m)
				removeState(m)
			}
		}(m)
	}
}

func addState(m Message) {
	fields := strings.Fields(m.Text)
	beginTime := time.Now().Unix()
	worklog := WorkLog{User: m.User, BeginTime: beginTime, Task: fields[3]}
	mutex.Lock()
	state[m.User] = worklog
	mutex.Unlock()
}

func removeState(m Message) {
	mutex.Lock()
	delete(state, m.User)
	mutex.Unlock()
}
