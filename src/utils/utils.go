package utils

import (
	"fmt"
	"logger"

	"github.com/hpcloud/tail"
)

const TAG string = "UTILS"


func reader(file_path string, tail_chan chan string) {
	t, err := tail.TailFile(file_path, tail.Config{Follow: true})
	if err != nil {
		logger.LogError(fmt.Sprintf("Error following file %s: %s",
									file_path, err.Error()), TAG)
		return
	}

	for line := range t.Lines {
		tail_chan <- line.Text
	}
}


func Tail(file_path string) (chan string) {

	logger.LogInfo(fmt.Sprintf("Tail -f %s", file_path), TAG)
	var tail chan string = make(chan string, 100)
	go reader(file_path, tail)
	return tail
}
