package auth

import (
	"events"
	"fmt"
	"logger"
	"os"
	"runtime"
	"testing"
	"time"
)

const (
	log_dir string = "test_logs/auth"
	test_line string = "Jan 01 00:00:00 testUser sudo: pam_unix(sudo:auth): "+
					   "authentication failure"
)


func init() {
	var go_root string = os.Getenv("GOPATH")
	os.Mkdir(fmt.Sprintf("%s/test_logs", go_root), 0777)
	os.Mkdir(fmt.Sprintf("%s/%s", go_root, log_dir), 0777)
	var filename string
	_, filename, _, _ = runtime.Caller(0)
	var filepath string
	filepath = fmt.Sprintf("%s/%s/auth-%s.log", go_root, log_dir,
						   time.Now().Format("01-02-2006_03-04"))
	var rc int = logger.SetLogFile(filepath, "true")
	if rc != 0 {
		panic("Could not set logger log file!")
	}
	logger.LogTestStep(fmt.Sprintf("Testing: %s", filename))
}


func TestAuthParser(t *testing.T) {
	es := []struct {
		input string
		output events.Event 
	}{
		{test_line, events.Event{EventType:events.AuthFailure,
								 Date:"Jan-01-00:00:00",
								 User:"testUser"}},
	}

	for _, c := range es {
		logger.LogTestStep(
			"Create event channel and inject the test line to the parser")
		var events_chan chan events.Event = make(chan events.Event, 1)
		var lines chan string = make(chan string, 1)
		go func () {
			lines <-c.input
		}()
		var dummy_chan chan bool
		go AuthParser(lines, events_chan, dummy_chan)
		logger.LogTestStep("Verify a new event is generated")
		time.Sleep(time.Second * 1)
		select {
		case event, _ := <-events_chan:
			if event != c.output {
				t.Errorf("Event is different than the expected one: %v", event)
			}
		default:
			t.Errorf("No event was generated")
		}
	}
}
