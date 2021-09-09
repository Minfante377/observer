package events

import (
	"api"
	"context"
	"fmt"
	"logger"

	"google.golang.org/grpc"
)

const (
	tag string = "EVENTS"
	address string = ":8080"
	AuthFailure int = 0
	TopMemoryUser int = 1
	MemoryTh int = 2
	DiskUsageAlarm int = 3
	HiddenFileMod int = 4
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


func sendEvents(event_chan chan Event, fail chan bool) {
	logger.LogInfo("Setting up a connection to the server", tag)
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.LogError("Could not set up connection to the server", tag)
		fail <-true
		return
	}
	defer conn.Close()
	c := api.NewEventsClient(conn)
	ctx := context.Background()
	for true {
		event := <-event_chan
		_, err := c.NewEvent(ctx, &api.Event{EventType:int32(event.EventType),
							                 Date:event.Date,
											 User:event.User,
										     Pwd:event.Pwd,
									         Cmd:event.Cmd,
										     Pid:string(event.Pid),
									         Notes:event.Notes})
		if err != nil {
			logger.LogError(fmt.Sprintf(
				"Error sending event: %s", err.Error()), tag)
		}
	}
}


func EventHandler(events chan Event) {
	var send_events chan Event = make(chan Event, 10)
	var fail_con chan bool = make(chan bool, 1)
	go sendEvents(send_events, fail_con)
	for true {
		var event Event
		event = <-events
		select {
			case _ = <-fail_con:
				logger.LogInfo("Connection is down. "+
							   "Cannot send event to server", tag)
			default:
				send_events <-event
		}
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
		}else if event.EventType == HiddenFileMod {
			logger.LogInfo(fmt.Sprintf("[%s]: %s", event.Date, event.Notes),
									   tag)
		}else{
			logger.LogInfo(fmt.Sprintf("Event %d not implemented",
									   event.EventType), tag)
		}
	}
}
