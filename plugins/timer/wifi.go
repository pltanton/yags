package timer

import (
	"bytes"
	"io/ioutil"
	"math"
	"strconv"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

// NewWifi creates timer with displays a network connection signal
func NewWifi(name string) Timer {
	conf := viper.Sub("plugins." + name)
	setWifiDefaults(conf)
	iface := []byte(conf.GetString("interface"))
	task := func() string {
		var format string
		lvl := parseNetwork(iface)
		if lvl != -1 {
			format = conf.GetString("connected")
		} else {
			format = conf.GetString("disconnected")
		}
		return utils.ReplaceVar(format, "lvl", strconv.FormatFloat(math.Floor(lvl+.5), 'f', 0, 64))
	}
	timer := NewTimerFunc(name, task)
	timer.conf = conf
	return timer
}

func parseNetwork(iface []byte) float64 {
	dat, _ := ioutil.ReadFile("/proc/net/wireless")
	for _, line := range bytes.Split(dat, []byte("\n")) {
		if bytes.Index(bytes.TrimSpace(line), iface) == 0 {
			lvl, _ := strconv.ParseFloat(string(bytes.Fields(line)[3]), 64)
			lvl = lvl / 70 * -100
			return lvl
		}
	}
	return -1
}
