package filesystem

import (
	"events"
	"fmt"
	"logger"
	"os"
	"time"
	"utils"
)

const tag string = "FILESYSTEM"

var stop bool = false


func totalUsageHandler(events_chan chan events.Event, usage_th int) {
	var partitions []utils.Partition
	for true {
		if stop {
			logger.LogInfo("Stopping disk usage analysis", tag)
			return
		}
		partitions = utils.GetDiskUsage()
		for _, partition := range partitions {
			if partition.Used > usage_th {
				var event events.Event
				event.EventType = events.DiskUsageAlarm
				now := time.Now()
				event.Date = now.Format("Jan-02-15:04:05")
				event.Notes = fmt.Sprintf("Disk usage %d is over the %d Th.\n"+
				                          "FileSystem:%s\nMounted:%s",
										  partition.Used, usage_th,
										  partition.Fs, partition.Mount)
				events_chan <- event
			}
		}
		time.Sleep(time.Hour * 6)
	}
}


func checkFile(file string, events_chan chan events.Event, last_check int64) {
	f, err := os.Stat(file)
	if err != nil {
		return
	}
	mtime := f.ModTime()
	if mtime.Unix() > last_check {
		var event events.Event
		event.EventType = events.HiddenFileMod
		event.Date = mtime.Format("Jan-02-15:04:05")
		event.Notes = fmt.Sprintf("Hidden file %s created or modified", file)
		events_chan <- event
	}
}


func checkDir(path string, events_chan chan events.Event, last_check int64) {
	var dirs, files []string
	dirs, files = utils.GetFiles(path)
	for _, f := range files {
		go checkFile(f, events_chan, last_check)
	}
	for _, dir := range dirs {
		checkDir(dir, events_chan, last_check)
	}
}


func hiddenFilesHandler(events_chan chan events.Event) {
	var last_check int64
	var dirs, files []string
	for true {
		if stop {
			logger.LogInfo("Stopping hidden files analysis", tag)
			return
		}
		now := time.Now()
		last_check = now.Unix()
		dirs, files = utils.GetFiles("/")
		for _, f := range files {
			go checkFile(f, events_chan, last_check)
		}
		for _, dir := range dirs {
			go checkDir(dir, events_chan, last_check)
		}
		time.Sleep(time.Hour * 1)
	}
}


func FsHandler(events_chan chan events.Event, usage_th int) {
	stop = false
	go totalUsageHandler(events_chan , usage_th)
	go hiddenFilesHandler(events_chan)
}


func StopFsHandler() {
	stop = true
}
