package memory

import (
	"events"
	"fmt"
	"logger"
	"time"
	"utils"
)

const tag = "MEMORY"

var stop bool = false


func psMemoryHandler(events_chan chan events.Event) {
	var processes []utils.Process
	var top_process utils.Process
	top_process.Pid = -1
	for true {
		if stop {
			logger.LogInfo("Stopping memory usage analysis by ps", tag)
		}
		processes = utils.GetProcessesByMemory()
		if processes != nil {
			if processes[0].Pid != top_process.Pid {
				var event events.Event
				now := time.Now()
				event.Date = now.Format("Jan-02-15:04:05")
				event.EventType = events.TopMemoryUser
				event.User = processes[0].User
				event.Pid = processes[0].Pid
				event.Cmd = processes[0].Cmd
				top_process = processes[0]
				events_chan <- event
			}
		}else{
			logger.LogError("Error getting processes", tag)
		}
		time.Sleep(time.Minute * 1)
	}
}


func totalMemoryHandler(events_chan chan events.Event, memory_th float64) {
	var th_flag bool = false
	var total_memory, available_memory, memory_usage float64
	for true {
		if stop {
			logger.LogInfo("Stopping total memory usage analisys", tag)
			return
		}
		total_memory, available_memory = utils.GetTotalMemory()
		if total_memory > 0 {
			memory_usage = (total_memory - available_memory) / total_memory
			if !th_flag && memory_usage > memory_th {
				var event events.Event
				event.EventType = events.MemoryTh
				now := time.Now()
				event.Date = now.Format("Jan-02-15:04:05")
				event.Notes = fmt.Sprintf("Current: %f, Th: %f",
										  memory_usage * 100,
										  memory_th * 100)
				events_chan <- event
				th_flag = true
			} else if memory_usage < memory_th * 0.9 {
				th_flag = false
			}
		}

		time.Sleep(time.Minute * 1)
	}
}


func MemoryHandler(events_chan chan events.Event, memory_th float64) {
	logger.LogInfo("Starting memory usage analysis", tag)
	stop = false
	go psMemoryHandler(events_chan)
	go totalMemoryHandler(events_chan, memory_th)
}


func StopMemoryHandler() {
	stop = true
}
