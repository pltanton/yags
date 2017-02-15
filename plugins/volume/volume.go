package volume

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

func (v Volume) parseVolume() string {

	volumeStr, _ := exec.Command("pamixer", "--get-volume").Output()
	muteStr, _ := exec.Command("pamixer", "--get-mute").Output()

	volume, _ := strconv.Atoi(strings.TrimSpace(string(volumeStr)))
	muted, _ := strconv.ParseBool(strings.TrimSpace(string(muteStr)))

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
	ctl, err := openCtl("hw:0")
	defer closeCtl(ctl)
	if err != nil {
		panic(err)
	}

	v.out <- v.parseVolume()
	for {
		isEvent, err := pollCtl(ctl)
		if err != nil {
			v.out <- err.Error()
		} else if isEvent {
			v.out <- v.parseVolume()
		}
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
