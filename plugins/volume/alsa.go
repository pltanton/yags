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

func checkEvent(ctl *C.snd_ctl_t) (bool, error) {
	var event *C.snd_ctl_event_t

	C.snd_ctl_event_malloc(&event)
	defer C.snd_ctl_event_free(event)

	err := C.snd_ctl_read(ctl, event)
	if err < 0 {
		return false, errors.New("Cannot read event")
	}

	if C.snd_ctl_event_get_type(event) != C.SND_CTL_EVENT_ELEM {
		return false, nil
	}

	mask := C.snd_ctl_event_elem_get_mask(event)

	if (mask & C.SND_CTL_EVENT_MASK_VALUE) == 0 {
		return false, nil
	}

	return true, nil
}

func pollCtl(ctl *C.snd_ctl_t) (bool, error) {
	var fd C.struct_pollfd
	C.snd_ctl_poll_descriptors(ctl, &fd, 1)

	err := C.poll(&fd, 1, -1)
	if err <= 0 {
		return false, fmt.Errorf("Cannot poll: %v", err)
	}
	var revents C.ushort
	C.snd_ctl_poll_descriptors_revents(ctl, &fd, 1, &revents)
	if (revents & C.POLLIN) != 0 {
		return checkEvent(ctl)
	}
	return false, nil
}
