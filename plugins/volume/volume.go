// Created by cgo - DO NOT EDIT

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:1
package volume

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:10

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:9
import (
	"os/exec"
	"strconv"
	"strings"
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:15

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:14
	"github.com/spf13/viper"
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:17

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:16
	"github.com/pltanton/yags/utils"
)

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:20

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:19
func (v Volume) parseVolume() string {
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:22

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:21
	volumeStr, _ := exec.Command("pamixer", "--get-volume").Output()
	muteStr, _ := exec.Command("pamixer", "--get-mute").Output()
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:25

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:24
	volume, _ := strconv.Atoi(strings.TrimSpace(string(volumeStr)))
	muted, _ := strconv.ParseBool(strings.TrimSpace(string(muteStr)))
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:28

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:27
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
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:42

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:41
	return utils.ReplaceVar(pattern, "vol", strconv.Itoa(volume))
}

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:45

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:44
func (v Volume) StartMonitor() {
	ctl, err := openCtl("hw:0")
	defer func(_cgo0 *_Ctype_struct__snd_ctl) {
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:46
		_Cfunc_snd_ctl_close(_cgoCheckPointer((*_Ctype_struct__snd_ctl)(_cgo0)).(*_Ctype_struct__snd_ctl))
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:46
	}(ctl)
	if err != nil {
		panic(err)
	}
//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:52

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:51
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

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:63

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:62
type Volume struct {
	conf *viper.Viper
	out  chan string
}

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:68

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:67
func (v Volume) Chan() chan string {
	return v.out
}

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:73

//line /home/anton/go/src/github.com/pltanton/yags/plugins/volume/volume.go:72
func NewVolume(name string) Volume {
	return Volume{
		out:  make(chan string),
		conf: setDefaults(viper.Sub("plugins." + name)),
	}
}
