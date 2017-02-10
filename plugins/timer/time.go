package timer

import (
	"time"

	"github.com/spf13/viper"
)

// NewTime creates timer with a representing time task
func NewTime(name string) Timer {
	conf := viper.Sub("plugins." + name)
	setTimeDefaults(conf)
	task := func() string {
		format := conf.GetString("timeFormat")
		return time.Now().Format(format)
	}
	timer := NewTimerFunc(name, task)
	timer.conf = conf
	return timer
}
