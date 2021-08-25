package events

import (
	"fmt"
	"logger"
)

const (
	tag string = "EVENTS"
	AuthFailure int = 0
)

type Event struct {
	EventType int
	Date string
	User string
	Pwd string
	Cmd string
	Notes string
}


func EventHandler(events chan Event) {
	for true {
		var event Event
		event = <-events
		if event.EventType == AuthFailure {
			logger.LogInfo(fmt.Sprintf("Auth failure event:\n[%s]: user=%s",
						               event.Date,
							           event.User), tag)
		}else{
			logger.LogInfo(fmt.Sprintf("Event %d not implemented",
									   event.EventType), tag)
		}
	}
}
