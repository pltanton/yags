package volume

/*
#cgo LDFLAGS: -lasound
#include <alsa/asoundlib.h>
*/
import "C"

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

func (v Volume) parseVolume() string {
	// TODO: fetch volume, using native library
	volumeStr, _ := exec.Command("pamixer", "--get-volume").Output()
	muteStr, _ := exec.Command("pamixer", "--get-mute").Output()

	volume, _ := strconv.Atoi(strings.TrimSpace(string(volumeStr)))
	muted, _ := strconv.ParseBool(strings.TrimSpace(string(muteStr)))

	// TODO: set defalut values
	config := viper.Sub("plugins." + v.name)

	var pattern string
	if muted {
		pattern = config.GetString("muted")
	} else {
		switch {
		case volume > 66:
			pattern = config.GetString("hight")
		case volume > 33:
			pattern = config.GetString("medium")
		default:
			pattern = config.GetString("low")
		}
	}

	return utils.ReplaceVar(pattern, "vol", strconv.Itoa(volume))
}

func (v Volume) StartMonitor() {
	ctl, err := openCtl("hw:0")
	defer C.snd_ctl_close(ctl)
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
	name string
	out  chan string
}

func (v Volume) Chan() chan string {
	return v.out
}

// NewVolume returns an instance of volme
func NewVolume(name string) Volume {
	return Volume{
		name: name,
		out:  make(chan string),
	}
}
