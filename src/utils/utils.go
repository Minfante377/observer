package utils

import (
	"bufio"
	"fmt"
	"logger"
	"os"
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

type Partition struct {
	Used int
	Fs string
	Mount string
}

const tag string = "UTILS"


func reader(file_path string, tail_chan TailChannel) {
	var config tail.Config = tail.Config{Follow: true, MustExist: true}
	t, err := tail.TailFile(file_path, config)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error following file %s: %s",
									file_path, err.Error()), tag)
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


func ParsePsOutput(output string) []Process {
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


func ParseDfOutput(output string) []Partition {
	var partitions []Partition
	space := regexp.MustCompile(`\s+`)
	output = strings.TrimSuffix(output, "\n")
	for _, line := range strings.Split(output, "\n") {
		line = space.ReplaceAllString(line, " ")
		elements := strings.Split(line, " ")
		var partition Partition
		partition.Fs = elements[0]
		partition.Mount = elements[5]
		partition.Used, _ = strconv.Atoi(strings.Replace(
			elements[4], "%", "", 1))
		partitions = append(partitions, partition)
	}
	return partitions
}


func Tail(file_path string) (TailChannel) {

	logger.LogInfo(fmt.Sprintf("Tail -f %s", file_path), tag)
	var tail TailChannel
	tail.Tail = make(chan string, 100)
    tail.Stop = make(chan bool, 1)
	go reader(file_path, tail)
	return tail
}


func GetTotalMemory() (float64, float64) {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		logger.LogError("Error reading meminfo file", tag)
		return -1 , -1
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	var total_memory_line string = scanner.Text()
	scanner.Scan()
	scanner.Scan()
	var available_memory_line string = scanner.Text()
	space := regexp.MustCompile(`\s+`)
	total_memory_line = space.ReplaceAllString(total_memory_line, " ")
	available_memory_line = space.ReplaceAllString(available_memory_line, " ")
	total_memory_line = strings.Split(total_memory_line, " ")[1]
	available_memory_line = strings.Split(available_memory_line, " ")[1]
	total_memory, _ := strconv.ParseFloat(total_memory_line, 32)
	available_memory, _ := strconv.ParseFloat(available_memory_line, 32)
	return total_memory, available_memory
}


func GetProcessesByMemory() []Process {
	var cmd string = "ps aux --sort -rss | head -n 6| tail -n 5"
	command := exec.Command("bash", "-c", cmd)
	out, err := command.Output()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error getting processes by memory: %s",
									err.Error()), tag)
		return nil
	}
	processes := ParsePsOutput(string(out))
	return processes
}


func GetDiskUsage() []Partition {
	var cmd string = "df -h --type=ext4| tail -n +2"
	command := exec.Command("bash", "-c", cmd)
	out, err := command.Output()
	if err != nil {
		logger.LogError(fmt.Sprintf("Error getting disk usage: %s",
									err.Error()), tag)
		return nil
	}
	partitions := ParseDfOutput(string(out))
	return partitions
}
