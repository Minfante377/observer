package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"logger"
	"os"
	"os/exec"
	"path/filepath"
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

type Config struct {
	Auth bool
	Memory bool
	Fs bool
	MemoryTh float64
	StorageTh int
}


const tag string = "UTILS"
const AuthModule string = "auth_module"
const MemoryModule string = "memory_module"
const MemoryTh string = "memory_th"
const FsModule string = "filesystem_module"
const FsTh string = "filesystem_th"


func reader(file_path string, tail_chan TailChannel) {
	var config tail.Config = tail.Config{Follow: true, MustExist: true}
	t, err := tail.TailFile(file_path, config)
	if err != nil {
		logger.LogError(fmt.Sprintf("Error following file %s: %s",
									file_path, err.Error()), tag)
		return
	}

	for true {
		line := <-t.Lines
		select {
			case <-tail_chan.Stop:
				return
			default:
				tail_chan.Tail <- line.Text
		}
	}
}


func parseConfigLine (line string) string {
	line = strings.Split(line, "=")[1]
	line = strings.TrimSpace(line)
	return line
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


func GetFiles(path string) ([]string, []string) {
	files, err := ioutil.ReadDir(path)
	var dirs, hidden_files []string
	if err != nil {
		return nil, nil
	}

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			dirs = append(dirs, filepath.Join(path, string(fileInfo.Name())))
		}else{
			if string(fileInfo.Name()[0]) == "." {
				hidden_files = append(
					hidden_files,
					filepath.Join(path, string(fileInfo.Name())))
			}
		}
	}
	return dirs, hidden_files
}


func ReadConfig(path string) Config {
	f, err := os.Open(path)
	var config Config
	config.Auth = false
	config.Memory = false
	config.Fs = false
	config.MemoryTh = 1.0
	config.StorageTh = 100
	if err != nil {
		logger.LogError(fmt.Sprintf("Could not open config file on %s: %s",
									path, err.Error()), tag)
		return config
	}
	defer f.Close()

	logger.LogInfo(fmt.Sprintf("Parsing config from: %s", path), tag)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var line string
		line = scanner.Text()
		if strings.Contains(line, AuthModule) {
			auth_module, err := strconv.Atoi(parseConfigLine(line))
			if err == nil {
				if auth_module == 1 {
					config.Auth = true
				}
			}else {
				logger.LogError(fmt.Sprintf("Error parsing %s", AuthModule),
								tag)
			}
		}else if strings.Contains(line, MemoryModule) {
			memory_module, err := strconv.Atoi(parseConfigLine(line))
			if err == nil {
				if memory_module == 1 {
					config.Memory = true
				}
			}else {
				logger.LogError(fmt.Sprintf("Error parsing %s", MemoryModule),
								tag)
			}
		}else if strings.Contains(line, FsModule) {
			fs_module, err := strconv.Atoi(parseConfigLine(line))
			if err == nil {
				if fs_module == 1 {
					config.Fs = true
				}
			}else {
				logger.LogError(fmt.Sprintf("Error parsing %s", FsModule), tag)
			}
		}else if strings.Contains(line, MemoryTh) {
			memory_th, err := strconv.ParseFloat(parseConfigLine(line), 32)
			if err == nil {
				config.MemoryTh = memory_th
			}else {
				logger.LogError(fmt.Sprintf("Error parsing %s", MemoryTh), tag)
			}
		}else if strings.Contains(line, FsTh) {
			fs_th, err := strconv.ParseFloat(parseConfigLine(line), 32)
			if err == nil {
				config.StorageTh = int(fs_th * 100)
			}else {
				logger.LogError(fmt.Sprintf("Error parsing %s", FsTh), tag)
			}
		}
	}
	logger.LogInfo(fmt.Sprintf("Parsed config is:\nAuth=%t\nMemory=%t\n"+
							   "Memory TH=%f\nFs=%t\nFs TH=%d", config.Auth,
							   config.Memory, config.MemoryTh, config.Fs,
							   config.StorageTh), tag)
	return config
}
