package timer

import (
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Producer type is a function, which should return stirng
type Producer func() string

// Timer just a struct to handl timer plugin
type Timer struct {
	task Producer
	out  chan string
	conf *viper.Viper
}

// Chan returns timer channel
func (t Timer) Chan() chan string { return t.out }

// StartMonitor returns the task execution result with a configured pause
// interval
func (t Timer) StartMonitor() {
	for {
		pause := t.conf.GetInt64("pause")
		t.out <- t.task()
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
}

// NewTimerCMD creates timer plugin with function builded from configured param
func NewTimerCMD(name string) Timer {
	return NewTimerFunc(
		name,
		func() string {
			cmd := viper.GetString("plugins." + name + ".cmd")
			res, err := exec.Command(cmd).Output()
			if err != nil {
				return err.Error()
			}
			return strings.TrimSpace(string(res))
		},
	)
}

// NewTimerFunc creates timer plugin with given function
func NewTimerFunc(name string, task Producer) Timer {
	return Timer{
		out:  make(chan string),
		task: task,
		conf: viper.Sub("plugins." + name),
	}
}
