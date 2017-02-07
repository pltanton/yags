package volume

/*
#cgo LDFLAGS: -lasound
#include <alsa/asoundlib.h>
#include <stdio.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"os/exec"
)

func openCtl(name string) (*C.snd_ctl_t, error) {
	deviceCString := C.CString(name)
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
	volume, _ := exec.Command("pamixer", "--get-volume").Output()
	mute, _ := exec.Command("pamixer", "--get-mute").Output()
	fmt.Println(string(volume), string(mute))
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
