package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

var state map[string]WorkLog
var mutex *sync.Mutex
var worklogs []WorkLog
var cnFacts []string

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: hibot token\n")
		os.Exit(1)
	}

	// start a cancel channel
	cancel := make(chan Message)
	state = make(map[string]WorkLog)
	worklogs = make([]WorkLog, 1)
	cnFacts = getCnFacts()
	mutex = &sync.Mutex{}
	ws, id := slackConnect(os.Args[1])
	fmt.Println("hibot is ready to rock, ^C")

	go ping(ws)

	for {
		// read each incoming message
		m, err := getMessage(ws)
		if err != nil {
			log.Fatal(err)
		}

		handleMessage(ws, id, m, cancel)
	}
}
