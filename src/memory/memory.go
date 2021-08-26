package memory

import(
	"events"
	"logger"
	"time"
	"utils"
)

const tag = "MEMORY"

var stop bool = false


func MemoryHandler(events_chan chan events.Event) {
	var processes []utils.Process
	var top_process utils.Process
	top_process.Pid = -1
	logger.LogInfo("Starting memory usage analysis", tag)
	for true {
		if stop {
			logger.LogInfo("Stopping memory usage analysis", tag)
			stop = false
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


func StopMemoryHandler() {
	stop = true
}
