package main

import(
	"auth"
	"events"
	"filesystem"
	"logger"
	"memory"
	"time"
)

const tag string = "MAIN"
const memoryTh float64 = 0.8
const diskTh int = 10


func main() {
	var events_chan chan events.Event = make(chan events.Event, 100)
	go events.EventHandler(events_chan)
	go auth.AuthHandler(events_chan)
	go memory.MemoryHandler(events_chan, memoryTh)
	go filesystem.FsHandler(events_chan, diskTh)
	for true{
		logger.LogInfo("Running...", tag)
		time.Sleep(time.Hour * 1)
	}
}
