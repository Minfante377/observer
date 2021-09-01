package main

import(
	"auth"
	"events"
	"filesystem"
	"logger"
	"memory"
	"time"
	"utils"
)

const tag string = "MAIN"
const configPath string = "modules.conf"


func main() {
	var events_chan chan events.Event = make(chan events.Event, 100)
	var config utils.Config = utils.ReadConfig(configPath)
	go events.EventHandler(events_chan)
	if config.Auth {
		go auth.AuthHandler(events_chan)
	}
	if config.Memory {
		go memory.MemoryHandler(events_chan, config.MemoryTh)
	}
	if config.Fs {
		go filesystem.FsHandler(events_chan, config.StorageTh)
	}
	for true{
		logger.LogInfo("Running...", tag)
		time.Sleep(time.Hour * 1)
	}
}
