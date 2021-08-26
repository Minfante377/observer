package utils

import (
	"fmt"
	"logger"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/tail"
)

type TailChannel struct {
	Tail chan string
	Stop chan bool
}

type Process struct {
	User string
	Pid int
	Memory float64
	Cpu float64
	Cmd string
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


func parsePsOutput(output string) []Process {
	var processes []Process
	space := regexp.MustCompile(`\s+`)
	output = strings.TrimSuffix(output, "\n")
	for _, line := range strings.Split(output, "\n") {
		line = space.ReplaceAllString(line, " ")
		elements := strings.Split(line, " ")
		var process Process
		process.User = elements[0]
		process.Pid, _ = strconv.Atoi(elements[1])
		process.Cpu, _ = strconv.ParseFloat(elements[2], 32)
		process.Memory, _ = strconv.ParseFloat(elements[3], 32)
		process.Cmd = strings.Join(elements[10:], " ")
		processes = append(processes, process)
	}
	return processes
}

func Tail(file_path string) (TailChannel) {

	logger.LogInfo(fmt.Sprintf("Tail -f %s", file_path), TAG)
	var tail TailChannel
	tail.Tail = make(chan string, 100)
    tail.Stop = make(chan bool, 1)
	go reader(file_path, tail)
	return tail
}


func GetProcessesByMemory() []Process {
	var cmd string = "ps aux --sort -rss | head -n 6| tail -n 5"
	command := exec.Command("bash", "-c", cmd)
	out, err := command.Output()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error getting processes by memory: %s",
									err.Error()), TAG)
		return nil
	}
	processes := parsePsOutput(string(out))
	return processes
}
