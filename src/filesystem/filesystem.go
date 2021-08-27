package filesystem

import (
	"events"
	"fmt"
	"logger"
	"time"
	"utils"
)

const tag string = "FILESYSTEM"

var stop bool = false


func totalUsageHandler(events_chan chan events.Event, usage_th int) {
	var partitions []utils.Partition
	for true{
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


func FsHandler(events_chan chan events.Event, usage_th int) {
	stop = false
	go totalUsageHandler(events_chan , usage_th)
}


func StopFsHandler() {
	stop = true
}
