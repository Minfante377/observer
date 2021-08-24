package utils

import (
	"fmt"
	"logger"

	"github.com/hpcloud/tail"
)

type TailChannel struct {
	tail chan string
	stop chan bool
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
		if len(tail_chan.stop) == 0 {
			tail_chan.tail <- line.Text
		}else{
			return
		}
	}
}


func Tail(file_path string) (chan string, chan bool) {

	logger.LogInfo(fmt.Sprintf("Tail -f %s", file_path), TAG)
	var tail TailChannel
	tail.tail = make(chan string, 100)
	tail.stop = make(chan bool, 1)
	go reader(file_path, tail)
	return tail.tail, tail.stop
}
