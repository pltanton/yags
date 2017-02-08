package volume

/*
#cgo LDFLAGS: -lasound
#include <alsa/asoundlib.h>
#include <stdio.h>
*/
import "C"
import "unsafe"

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/pltanton/yags/utils"
	"github.com/spf13/viper"
)

func openCtl(name string) (*C.snd_ctl_t, error) {
	deviceCString := C.CString(name)
	defer C.free(unsafe.Pointer(deviceCString))
	var ctl *C.snd_ctl_t
	err := C.snd_ctl_open(&ctl, deviceCString, C.SND_CTL_READONLY)
	if err < 0 {
		msg := fmt.Sprintf("Cannot open ctl '%v', err: %v\n", name, err)
		return nil, errors.New(msg)
	}
	err = C.snd_ctl_subscribe_events(ctl, 1)
	if err < 0 {
		msg := fmt.Sprintf(
			"Cannot open subscribe events to ctl '%v', err: %v\n", name, err,
		)
		C.snd_ctl_close(ctl)
		return nil, errors.New(msg)
	}
	return ctl, nil
}

func checkEvent(ctl *C.snd_ctl_t) error {
	var event *C.snd_ctl_event_t

	C.snd_ctl_event_malloc(&event)
	defer C.free(unsafe.Pointer(event))

	err := C.snd_ctl_read(ctl, event)
	if err < 0 {
		return errors.New("Cannot read event")
	}

	if C.snd_ctl_event_get_type(event) != C.SND_CTL_EVENT_ELEM {
		return nil
	}

	mask := C.snd_ctl_event_elem_get_mask(event)

	if (mask & C.SND_CTL_EVENT_MASK_VALUE) == 0 {
		return nil
	}

	parseVolume()
	return nil
}

func parseVolume() {
	// TODO: fetch volume, using native library
	volumeStr, _ := exec.Command("pamixer", "--get-volume").Output()
	muteStr, _ := exec.Command("pamixer", "--get-mute").Output()

	volume, _ := strconv.Atoi(strings.TrimSpace(string(volumeStr)))
	muted, _ := strconv.ParseBool(strings.TrimSpace(string(muteStr)))

	// TODO: set default values
	config := viper.Sub("plugins.volume")

	var pattern string
	if muted {
		pattern = config.GetString("muted")
	} else {
		switch {
		case volume > 66:
			pattern = config.GetString("high")
		case volume > 33:
			pattern = config.GetString("medium")
		default:
			pattern = config.GetString("low")
		}
	}
	fmt.Println(utils.ReplaceVar(pattern, "vol", strconv.Itoa(volume)))
}

func Monitor() {
	ctl, err := openCtl("hw:0")
	defer C.snd_ctl_close(ctl)
	if err != nil {
		panic(err)
	}
	for {
		var fd C.struct_pollfd
		C.snd_ctl_poll_descriptors(ctl, &fd, 1)

		err := C.poll(&fd, 1, -1)
		if err <= 0 {
			break
		}
		var revents C.ushort
		C.snd_ctl_poll_descriptors_revents(ctl, &fd, 1, &revents)
		if (revents & C.POLLIN) != 0 {
			checkEvent(ctl)
		}
	}
}
