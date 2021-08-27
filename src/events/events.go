package events

import (
	"fmt"
	"logger"
)

const (
	tag string = "EVENTS"
	AuthFailure int = 0
	TopMemoryUser int = 1
	MemoryTh int = 2
	DiskUsageAlarm int = 3
)

type Event struct {
	EventType int
	Date string
	User string
	Pwd string
	Cmd string
	Pid int
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
		}else if event.EventType == TopMemoryUser {
			logger.LogInfo(
				fmt.Sprintf(
					"Top memory user changed:\n[%s] user=%s, pid=%d, cmd=%s",
					event.Date, event.User, event.Pid, event.Cmd), tag)
		}else if event.EventType == MemoryTh {
			logger.LogInfo(fmt.Sprintf("Memory th excedeed: %s", event.Notes),
									   tag)
		}else if event.EventType == DiskUsageAlarm {
			logger.LogInfo(fmt.Sprintf("Disk usage alarm: %s", event.Notes),
									   tag)
		}else{
			logger.LogInfo(fmt.Sprintf("Event %d not implemented",
									   event.EventType), tag)
		}
	}
}
