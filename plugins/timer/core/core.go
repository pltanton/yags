package core

import (
	"time"

	"github.com/spf13/viper"
)

// Producer type is a function, which should return stirng
type Producer func() string

// Timer just a struct to handle timer plugin
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

// NewTimerFunc creates timer plugin with given function
func NewTimerFunc(conf *viper.Viper, task Producer) Timer {
	return Timer{
		out:  make(chan string, 1),
		task: task,
		conf: conf,
	}
}
