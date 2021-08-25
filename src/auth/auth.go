package auth

import (
	"events"
	"fmt"
	"logger"
	"strings"
	"utils"
)


const (
	tag string = "AUTH.LOG"
	auth_file string = "/var/log/auth.log"
	auth_failure string = "authentication failure"
)

var tail_channel utils.TailChannel


func authParser(lines chan string, events_chan chan events.Event) {
	var line string
	for true {
		line = <-lines
		var fields []string = strings.Fields(line)
		var date string = strings.Join(fields[0:3], "-")
		var user string = fields[3]
		var event events.Event
		event.EventType = events.AuthFailure
		event.Date = date
		event.User = user
		events_chan <- event
	}
}


func isAuthEvent(line string) bool {
	var isAuth bool = strings.Contains(line, auth_failure)
	return isAuth
}


func AuthHandler(events_chan chan events.Event) {
	tail_channel = utils.Tail(auth_file)
	var auth_chan chan string = make(chan string, 10)
	go authParser(auth_chan, events_chan)
	for true{
		select {
			case <-tail_channel.Stop:
				logger.LogInfo(fmt.Sprintf("%s analysis stopped", auth_file),
							   tag)
				return
			case line := <-tail_channel.Tail:
				if isAuthEvent(line) {
					auth_chan <-line
				}
		}
	}
}


func StopAuthHandler() {
	tail_channel.Stop <- true
}
