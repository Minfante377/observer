package utils

import (
	"fmt"
	"logger"
	"os"
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
		tail, stop := Tail(c.input)

		logger.LogTestStep(
			"Write a line to the file and verify it was sent through the chan")
		_, err = test_file.WriteString("This is a test\n")
		if err != nil {
			t.Errorf("Error writing to the test file")
		}
		line, ok := <-tail
		if !ok {
			t.Errorf("Error on the channel")
		}
		if line != "This is a test" {
			t.Errorf("Information received is wrong %s vs This is a test",
			         line)
		}

		logger.LogTestStep("Stop the tail and verify the chan is closed")
		go func() {
			stop <- true
		}()
		time.Sleep(time.Second * 1)
		select {
		case _, ok = <-tail:
			t.Errorf("Channel did not get closed")
		default:
		}
	}
}
