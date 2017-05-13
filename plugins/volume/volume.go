package main

import (
	"strconv"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/plugins"
	"github.com/pltanton/yags/utils"
)

func (v volume) formatVolume(volume int, muted bool) string {
	var pattern string
	if muted {
		pattern = v.conf.GetString("muted")
	} else {
		switch {
		case volume > 66:
			pattern = v.conf.GetString("high")
		case volume > 33:
			pattern = v.conf.GetString("medium")
		default:
			pattern = v.conf.GetString("low")
		}
	}

	return utils.ReplaceVar(pattern, "vol", strconv.Itoa(volume))
}

func (v volume) StartMonitor() {
	client := newPulseClient()

	v.out <- v.formatVolume(client.getVolume(v.conf.GetString("sink")))
	for range client.event {
		v.out <- v.formatVolume(client.getVolume(v.conf.GetString("sink")))
	}
}

type volume struct {
	conf *viper.Viper
	out  chan string
}

func (v volume) Chan() chan string {
	return v.out
}

func New(conf *viper.Viper) plugins.Plugin {
	return volume{
		out:  make(chan string, 1),
		conf: conf,
	}
}
