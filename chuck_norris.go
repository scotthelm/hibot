package main

import (
	"bufio"
	"math/rand"
	"os"
	"time"

	"golang.org/x/net/websocket"
)

func getCnFacts() []string {
	cnFacts, err := readLines("./chuck_norris_facts.txt")
	if err != nil {
		cnFacts = []string{"No Chuck Norris Facts. Sad."}
	}
	return cnFacts
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func doChuckNorris(ws *websocket.Conn, m Message) {
	rand.Seed(time.Now().Unix())
	m.Text = cnFacts[rand.Intn(len(cnFacts))]
	postMessage(ws, m)
}
