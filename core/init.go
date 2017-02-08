package core

import (
	"fmt"

	"github.com/pltanton/yags/plugins"
	"github.com/pltanton/yags/plugins/volume"
)

var monitors map[string]plugins.Plugin

func Init() {
	monitors = make(map[string]plugins.Plugin)
	vol := volume.NewVolume("volume")
	monitors["volume"] = vol
	go vol.StartMonitor()
	for {
		msg := <-vol.Chan()
		fmt.Println(msg)
	}
	fmt.Println("end")
}
