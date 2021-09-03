package main

import (
	"auth"
	"events"
	"filesystem"
	"fmt"
	"logger"
	"memory"
	"os"
	"time"
	"utils"
)

const tag string = "MAIN"
const configPath string = "modules.conf"
var LogDir string
var Debug string


func initModules(config utils.Config, events_chan chan events.Event) {
	if config.Auth {
		go auth.AuthHandler(events_chan)
	}
	if config.Memory {
		go memory.MemoryHandler(events_chan, config.MemoryTh)
	}
	if config.Fs {
		go filesystem.FsHandler(events_chan, config.StorageTh)
	}
}


func stopModules() {
	go auth.StopAuthHandler()
	go memory.StopMemoryHandler()
	go filesystem.StopFsHandler()
}


func configHandler(events_chan chan events.Event) {
	now := time.Now()
	last_check := now.Unix()
	for true {
		f, err := os.Stat(configPath)
		if err != nil {
			logger.LogError(fmt.Sprintf("Error checking config file: %s",
										err.Error()), tag)
			continue
		}
		mod_time := f.ModTime().Unix()
		if mod_time > last_check {
			last_check = mod_time
			var config utils.Config = utils.ReadConfig(configPath)
			stopModules()
			initModules(config, events_chan)
		}
		time.Sleep(time.Second * 10)
	}
}


func main() {
	logger.SetLogFile(LogDir, Debug)
	var events_chan chan events.Event = make(chan events.Event, 100)
	var config utils.Config = utils.ReadConfig(configPath)
	go events.EventHandler(events_chan)
	initModules(config, events_chan)
	go configHandler(events_chan)
	for true{
		logger.LogInfo("Running...", tag)
		time.Sleep(time.Hour * 1)
	}
}
