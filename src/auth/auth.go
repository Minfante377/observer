package auth

import (
	"events"
	"fmt"
	"logger"
	"strings"
	"sync"
	"utils"
)


const (
	tag string = "AUTH.LOG"
	auth_file string = "/var/log/auth.log"
	auth_failure string = "authentication failure"
)

var tail_channels []utils.TailChannel
var stop_chans []chan bool
var mu sync.Mutex


func AuthParser(lines chan string, events_chan chan events.Event,
				stop chan bool) {
	var line string
	for true {
		select{
		case <-stop:
			return
		default:
		}
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
	var tail_channel utils.TailChannel = utils.Tail(auth_file)
	var auth_chan chan string = make(chan string, 10)
	var stop chan bool
	mu.Lock()
	tail_channels = append(tail_channels, tail_channel)
	stop_chans = append(stop_chans, stop)
	mu.Unlock()
	go AuthParser(auth_chan, events_chan, stop)
	logger.LogInfo("Starting auth.log analisys", tag)
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
	mu.Lock()
	if len(tail_channels) > 0 {
		tail_channels[0].Stop <- true
		stop_chans[0] <-true
		tail_channels = tail_channels[1:]
		stop_chans = stop_chans[1:]
	}
	mu.Unlock()
}
