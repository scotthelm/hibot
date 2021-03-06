package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/websocket"
)

func handleMessage(ws *websocket.Conn, id string, m Message, cancel chan Message) {
	// see if we're mentioned
	if m.Type == "message" && strings.HasPrefix(m.Text, "<@"+id+">") {
		// if so try to parse if
		parts := strings.Fields(m.Text)
		if len(parts) == 3 && parts[1] == "stock" {
			// looks good, get the quote and reply with the result
			go func(m Message) {
				m.Text = getQuote(parts[2])
				postMessage(ws, m)
			}(m)
			// NOTE: the Message object is copied, this is intentional
		} else if len(parts) == 5 && parts[1] == "timer" && parts[2] == "start" {
			go func(m Message) {
				doTimer(ws, m, cancel)
			}(m)
		} else if len(parts) == 3 && parts[1] == "timer" && parts[2] == "stop" {
			cancel <- m
		} else if len(parts) == 3 && parts[1] == "timer" && parts[2] == "report" {
			go func(m Message) {
				timerReport(ws, m)
			}(m)
		} else if len(parts) >= 2 && parts[1] == "chuck" {
			go func(m Message) {
				doChuckNorris(ws, m)
			}(m)
		} else {
			// huh?
			m.Text = fmt.Sprintf("sorry, that does not compute\n")
			postMessage(ws, m)
		}
	}

}
