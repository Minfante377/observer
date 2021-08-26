package memory

import(
	"events"
	"fmt"
	"logger"
	"time"
	"utils"
)

const tag = "MEMORY"

var stop bool = false


func MemoryHandler(events_chan chan events.Event) {
	var processes []utils.Process
	logger.LogInfo("Starting memory usage analysis", tag)
	for true {
		if stop {
			logger.LogInfo("Stopping memory usage analysis", tag)
			stop = false
		}
		processes = utils.GetProcessesByMemory()
		if processes != nil {
			logger.LogInfo(fmt.Sprintf("The process using more memory is: %v",
								   	   processes[0]), tag)
		}else{
			logger.LogError("Error getting processes", tag)
		}
		time.Sleep(time.Minute * 1)
	}
}


func StopMemoryHandler() {
	stop = true
}
