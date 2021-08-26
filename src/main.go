package main

import(
	"auth"
	"events"
	"logger"
	"memory"
	"time"
)

const tag string = "MAIN"


func main() {
	var events_chan chan events.Event = make(chan events.Event, 100)
	go events.EventHandler(events_chan)
	go auth.AuthHandler(events_chan)
	go memory.MemoryHandler(events_chan)
	for true{
		logger.LogInfo("Running...", tag)
		time.Sleep(time.Hour * 1)
	}
}
