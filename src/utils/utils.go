package utils

import (
	"fmt"
	"logger"

	"github.com/hpcloud/tail"
)

type TailChannel struct {
	Tail chan string
	Stop chan bool
}

const TAG string = "UTILS"


func reader(file_path string, tail_chan TailChannel) {
	var config tail.Config = tail.Config{Follow: true, MustExist: true}
	t, err := tail.TailFile(file_path, config)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error following file %s: %s",
									file_path, err.Error()), TAG)
		return
	}

	for line := range t.Lines {
		select {
			case <-tail_chan.Stop:
				return
			default:
				tail_chan.Tail <- line.Text
		}
	}
}


func Tail(file_path string) (TailChannel) {

	logger.LogInfo(fmt.Sprintf("Tail -f %s", file_path), TAG)
	var tail TailChannel
	tail.Tail = make(chan string, 100)
	tail.Stop = make(chan bool, 1)
	go reader(file_path, tail)
	return tail
}
