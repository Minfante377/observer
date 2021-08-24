package main

import(
	"fmt"
	"logger"
	"utils"
)

const TAG string = "MAIN"


func main() {
	var sys_tail chan string
	sys_tail = utils.Tail("/var/log/auth.log")
	for true{
		var line string = <-sys_tail
		logger.LogInfo(fmt.Sprintf("Last line is: %s", line), TAG)
	}
}
