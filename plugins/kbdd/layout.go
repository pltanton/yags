package kbdd

import (
	"fmt"

	"github.com/godbus/dbus"
	"github.com/spf13/viper"
)

// KBDD plugin structure
type KBDD struct {
	conf *viper.Viper
	out  chan string
}

// NewKBDD returns new instance of battery plugin by given name
func NewKBDD(name string) KBDD {
	return KBDD{
		out:  make(chan string, 1),
		conf: viper.Sub("plugins." + name),
	}
}

// Chan returns a strings channel with layout state monitor
func (k KBDD) Chan() chan string { return k.out }

// StartMonitor starts monitoring for battery changing events
func (k KBDD) StartMonitor() {
	conn, err := dbus.SessionBus()
	if err != nil {
		panic(fmt.Errorf("cannot connect dbus session: %s", err.Error()))
	}

	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		"path=/ru/gentoo/KbddService,interface=ru.gentoo.kbdd,"+
			"member=layoutChanged",
	)

	k.sendLayout(askForCurLayout())
	c := make(chan *dbus.Signal, 1)
	conn.Signal(c)
	for v := range c {
		k.sendLayout(v.Body[0].(uint32))
	}
}

func (k KBDD) sendLayout(layout uint32) {
	k.out <- k.layoutToString(layout)
}

func askForCurLayout() uint32 {
	conn, _ := dbus.SessionBus()
	obj := conn.Object("ru.gentoo.KbddService", "/ru/gentoo/KbddService")
	var layout uint32
	obj.Call("ru.gentoo.kbdd.getCurrentLayout", 0).Store(&layout)
	return layout
}

func (k KBDD) layoutToString(layout uint32) string {
	layouts := k.conf.GetStringSlice("names")
	return layouts[layout]
}
