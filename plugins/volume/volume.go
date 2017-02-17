package volume

import (
	"strconv"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

func (v Volume) formatVolume(volume int, muted bool) string {
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

func (v Volume) StartMonitor() {
	client := newPulseClient()

	v.out <- v.formatVolume(client.getVolume(v.conf.GetString("sink")))
	for range client.event {
		v.out <- v.formatVolume(client.getVolume(v.conf.GetString("sink")))
	}
}

type Volume struct {
	conf *viper.Viper
	out  chan string
}

func (v Volume) Chan() chan string {
	return v.out
}

func NewVolume(name string) Volume {
	return Volume{
		out:  make(chan string, 1),
		conf: setDefaults(viper.Sub("plugins." + name)),
	}
}
