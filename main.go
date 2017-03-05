package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: hibot token\n")
		os.Exit(1)
	}

	// start a cancel channel
	cancel := make(chan Message)
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
