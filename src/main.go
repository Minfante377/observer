package main

import(
	"fmt"
	"logger"
	"utils"
)

const TAG string = "MAIN"


func main() {
	var sys_tail chan string
	var stop chan bool
	var count int = 0
	sys_tail, stop = utils.Tail("/var/log/auth.log")
	for true{
		if count > 0 {
			stop <- true
			logger.LogInfo("Stopping tail", TAG)
			break
		}
		line, ok := <-sys_tail
		if ok {
			logger.LogInfo(fmt.Sprintf("Last line is: %s", line), TAG)
			count += 1
		} else {
			logger.LogError("Channel is closed", TAG)
			break
		}
	}
}
