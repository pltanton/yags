package timer

import (
	"time"

	"github.com/spf13/viper"
)

// NewTime creates timer with a representing time task
func NewTime(name string) Timer {
	task := func() string {
		format := viper.GetString("plugins." + name + ".timeFormat")
		return time.Now().Format(format)
	}
	return NewTimerFunc(name, task)
}
