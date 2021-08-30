package utils

import (
	"fmt"
	"logger"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)

const (
	log_dir string = "test_logs/utils"
	test_file string = "./test_file"
)


func init() {
	var go_root string = os.Getenv("GOPATH")
	os.Mkdir(fmt.Sprintf("%s/test_logs", go_root), 0777)
	os.Mkdir(fmt.Sprintf("%s/%s", go_root, log_dir), 0777)
	var filename string
	_, filename, _, _ = runtime.Caller(0)
	var filepath string
	filepath = fmt.Sprintf("%s/%s/utils-%s.log", go_root, log_dir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger.SetLogFile(filepath, "true")
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	logger.LogTestStep(fmt.Sprintf("Testing: %s", filename))
}


func TestTail(t *testing.T) {
	es := []struct {
		input          string
	}{
		{test_file},
	}

	logger.LogTestStep(
		"Create test file and start tailing it")
	for _, c := range es {
		logger.LogTestStep(fmt.Sprintf("Create %s", c.input))
		test_file, err := os.OpenFile(c.input,
									  os.O_APPEND|os.O_CREATE|os.O_WRONLY,
									  0666)
		if err != nil {
			t.Errorf("Error creating test file %s", c.input)
		}
		defer func() {
			logger.LogTestStep("Remove test files")
			os.RemoveAll(c.input)
		}()
		tail_channel:= Tail(c.input)

		logger.LogTestStep(
			"Write a line to the file and verify it was sent through the chan")
		_, err = test_file.WriteString("This is a test\n")
		if err != nil {
			t.Errorf("Error writing to the test file")
		}
		line, ok := <-tail_channel.Tail
		if !ok {
			t.Errorf("Error on the channel")
		}
		if line != "This is a test" {
			t.Errorf("Information received is wrong %s vs This is a test",
			         line)
		}

		logger.LogTestStep("Stop the tail and verify the chan is closed")
		go func() {
			tail_channel.Stop <- true
		}()
		time.Sleep(time.Second * 1)
		select {
		case _, ok = <-tail_channel.Tail:
			t.Errorf("Channel did not get closed")
		default:
		}
	}
}


func TestParsePsOutput(t *testing.T) {
	es := []struct {
		input          string
		output Process
	}{
		{"test_user 0  1.0 1.0 a a a   a   a a cmd --testarg",
		Process{User:"test_user", Pid:0, Cpu:1.0, Memory:1.0,
				Cmd:"cmd --testarg"}},
	}
	for _, c := range es {
		logger.LogTestStep(fmt.Sprintf(
			"Parse test input %s and verify the result", c.input))
		process := ParsePsOutput(c.input)[0]
		if process.User != c.output.User {
			t.Errorf("User was not parsed correctly: %s", process.User)
		}
		if process.Pid != c.output.Pid {
			t.Errorf("PID was not parsed correctly: %d", process.Pid)
		}
		if process.Cpu != c.output.Cpu {
			t.Errorf("CPU usage was not parsed correctly: %f", process.Cpu)
		}
		if process.Memory != c.output.Memory {
			t.Errorf("Memory was not parsed correctly: %f", process.Memory)
		}
		if process.Cmd != c.output.Cmd {
			t.Errorf("Cmd was not parsed correctly: %s", process.Cmd)
		}
	}
}


func TestParseDfOutput(t *testing.T) {
	es := []struct {
		input          string
		output Partition
	}{
		{"/dev/sda2       100G  50G  50G  50% /test/path",
		Partition{Fs:"/dev/sda2", Used:50, Mount:"/test/path"}},
	}
	for _, c := range es {
		logger.LogTestStep(fmt.Sprintf(
			"Parse test input %s and verify the result", c.input))
		partition := ParseDfOutput(c.input)[0]
		if partition.Fs != c.output.Fs {
			t.Errorf("Fs was not parsed correctly: %s", partition.Fs)
		}
		if partition.Used != c.output.Used {
			t.Errorf("Usage was not parsed correctly: %d", partition.Used)
		}
		if partition.Mount != c.output.Mount {
			t.Errorf("Mount path usage was not parsed correctly: %s",
					 partition.Mount)
		}
	}
}


func TestGetFiles(t *testing.T) {
	es := []struct {
		input          string
		output []string
	}{
		{"./test_folder",
		 []string{"test_folder", ".test_file"}},
	}
	for _, c := range es {
		logger.LogTestStep(fmt.Sprintf("Create test files"))
		os.Mkdir(c.input, 0777)
		os.Mkdir(filepath.Join(c.input, c.output[0]), 0777)
		_, err := os.OpenFile(filepath.Join(c.input, c.output[1]),
							  os.O_APPEND|os.O_CREATE|os.O_WRONLY,
							   0666)
		if err != nil {
			t.Errorf("Error creating test file %s", c.input)
		}
		defer func() {
			logger.LogTestStep("Remove test files")
			os.RemoveAll(c.input)
		}()

		logger.LogTestStep("Execute function and verify output")
		dirs, files := GetFiles(c.input)
		if dirs[0] != filepath.Join(c.input, c.output[0]) {
			t.Errorf("Failure retrieving folders: %s vs %s", dirs[0],
					 filepath.Join(c.input, c.output[0]))
		}
		if files[0] != filepath.Join(c.input, c.output[1]) {
			t.Errorf("Failure retrieving hidden files: %s vs %s", files[0],
			         filepath.Join(c.input, c.output[1]))
        }
	}
}
