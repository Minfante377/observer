package filesystem

import (
	"events"
	"fmt"
	"logger"
	"os"
	"sync"
	"time"
	"utils"
)

const tag string = "FILESYSTEM"

var stop_chans []chan bool
var mu sync.Mutex


func totalUsageHandler(events_chan chan events.Event, usage_th int,
					   stop chan bool) {
	var partitions []utils.Partition
	for true {
		select {
		case <-stop:
			logger.LogInfo("Stopping disk usage analysis", tag)
			return
		default:
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


func _checkDir(path string, events_chan chan events.Event, last_check int64) {
	var dirs, files []string
	dirs, files = utils.GetFiles(path)
	for _, f := range files {
		go checkFile(f, events_chan, last_check)
	}
	for _, dir := range dirs {
		_checkDir(dir, events_chan, last_check)
	}
}


func checkDir(path string, events_chan chan events.Event, last_check int64,
			  wg *sync.WaitGroup) {
	var dirs, files []string
	dirs, files = utils.GetFiles(path)
	for _, f := range files {
		go checkFile(f, events_chan, last_check)
	}
	for _, dir := range dirs {
		_checkDir(dir, events_chan, last_check)
	}
	wg.Done()
}


func hiddenFilesHandler(events_chan chan events.Event, stop chan bool) {
	var last_check int64
	var dirs, files []string
	var wg sync.WaitGroup
	now := time.Now()
	last_check = now.Unix()
	for true {
		select {
		case <-stop:
			logger.LogInfo("Stopping hidden files analysis", tag)
			return
		default:
		}
		dirs, files = utils.GetFiles("/")
		for _, f := range files {
			go checkFile(f, events_chan, last_check)
		}
		for _, dir := range dirs {
			wg.Add(1)
			go checkDir(dir, events_chan, last_check, &wg)
		}
		now = time.Now()
		last_check = now.Unix()
		wg.Wait()
		time.Sleep(time.Minute * 1)
	}
}


func FsHandler(events_chan chan events.Event, usage_th int) {
	var stop chan bool
	logger.LogInfo("Starting filesystem analysis", tag)
	mu.Lock()
	stop_chans = append(stop_chans, stop)
	mu.Unlock()
	go totalUsageHandler(events_chan , usage_th, stop)
	go hiddenFilesHandler(events_chan, stop)
}


func StopFsHandler() {
	mu.Lock()
	if len(stop_chans) > 0 {
		stop_chans[0] <- true
		stop_chans[0] <- true
		stop_chans = stop_chans[1:]
	}
	mu.Unlock()
}
