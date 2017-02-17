package volume

import (
	"github.com/godbus/dbus"
	"github.com/sqp/pulseaudio"
)

type pulseClient struct {
	*pulseaudio.Client
	event chan bool
}

func (pc *pulseClient) DeviceVolumeUpdated(path dbus.ObjectPath, values []uint32) {
	pc.event <- true
}

func (pc *pulseClient) DeviceMuteUpdated(path dbus.ObjectPath, values bool) {
	pc.event <- true
}

func (pc *pulseClient) getVolume(sink string) (vol int, mute bool) {
	dev := pc.Device(dbus.ObjectPath(sink))

	mute, _ = dev.Bool("Mute")
	vols, _ := dev.ListUint32("Volume")
	vol = int(volAvg(vols) * 100 / 65535)

	return
}

func volAvg(vols []uint32) (vol uint32) {
	if l := len(vols); l > 0 {
		for _, v := range vols {
			vol += v
		}
		vol /= uint32(l)
	}
	return vol
}

func newPulseClient() *pulseClient {
	pulse, err := pulseaudio.New()
	if err != nil {
		panic(err.Error())
	}

	client := &pulseClient{pulse, make(chan bool, 1)}
	pulse.Register(client)

	go pulse.Listen()

	return client
}
