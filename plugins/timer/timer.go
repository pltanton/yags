package timer

import (
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

// Producer type is a function, which should return stirng
type Producer func() string

// Timer just a struct to handl timer plugin
type Timer struct {
	name string
	task Producer
	out  chan string
}

// Chan returns timer channel
func (t Timer) Chan() chan string { return t.out }

// StartMonitor returns the task execution result with a configured pause
// interval
func (t Timer) StartMonitor() {
	for {
		pause := viper.GetInt64("plugins." + t.name + ".pause")
		format := viper.GetString("plugins." + t.name + ".format")
		t.out <- utils.ReplaceVar(format, "cmd", t.task())
		time.Sleep(time.Duration(pause) * time.Millisecond)
	}
}

func (t Timer) StopMonitor() {}

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
		name: name,
		out:  make(chan string),
		task: task,
	}
}
