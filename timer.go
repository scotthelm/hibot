package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type WorkLog struct {
	User      string
	BeginTime time.Time
	EndTime   time.Time
	Duration  time.Duration
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
		worklog = addState(m)
		minutes, err := minutes(m)
		if err != nil {
			m.Text = "Error creating task - make sure minutes is a number"
		} else {
			m.Text = fmt.Sprintf("Added a timer for %s", worklog.Task)
		}
		postMessage(ws, m)
		go func(m Message) {
			select {
			case _ = <-cancel:
				m.Text = fmt.Sprintf("<@%s>, your timer for %s has been canceled", worklog.User, worklog.Task)
				postMessage(ws, m)
				addLog(m)
				removeState(m)
			case _ = <-time.After(time.Duration(minutes) * time.Minute):
				m.Text = fmt.Sprintf("<@%s>, your time for %s has elapsed", worklog.User, worklog.Task)
				postMessage(ws, m)
				addLog(m)
				removeState(m)
			}
		}(m)
	}
}

func minutes(m Message) (int, error) {
	fields := strings.Fields(m.Text)
	return strconv.Atoi(fields[4])
}
func addState(m Message) WorkLog {
	fields := strings.Fields(m.Text)
	beginTime := time.Now()
	worklog := WorkLog{User: m.User, BeginTime: beginTime, Task: fields[3]}
	mutex.Lock()
	state[m.User] = worklog
	mutex.Unlock()
	return worklog
}

func removeState(m Message) {
	mutex.Lock()
	delete(state, m.User)
	mutex.Unlock()
}

func addLog(m Message) {
	mutex.Lock()
	worklog := state[m.User]
	worklog.EndTime = time.Now()
	worklogs = append(worklogs, worklog)
	mutex.Unlock()
}

func timerReport(ws *websocket.Conn, m Message) {
	logmessages := userLogMessages(m)
	if len(logmessages) > 0 {
		m.Text = strings.Join(logmessages[:], "\n")
	} else {
		m.Text = "You have no timed tasks."
	}
	postMessage(ws, m)
}

func userWorkLogs(m Message) []WorkLog {
	mutex.Lock()
	var userlogs []WorkLog = make([]WorkLog, 0)
	for _, wl := range worklogs {
		if wl.User == m.User {
			userlogs = append(userlogs, wl)
		}
	}
	mutex.Unlock()
	return userlogs
}

func userLogMessages(m Message) []string {
	var logmessages []string = make([]string, 0)
	for _, wl := range userWorkLogs(m) {
		logmessages = append(logmessages, wl.forSlack())
	}
	return logmessages
}

func (wl WorkLog) forSlack() string {
	wl.Duration = wl.EndTime.Sub(wl.BeginTime)
	return fmt.Sprintf("%s - %.2f", wl.Task, wl.Duration.Minutes())
}
