package battery

import (
	"fmt"
	"strconv"

	"github.com/godbus/dbus"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

// Battery plugin structure
type Battery struct {
	conf    *viper.Viper
	batName string
	out     chan string
}

// NewBattery returns new instance of battery plugin by given name
func NewBattery(name string) Battery {
	return Battery{
		out:  make(chan string),
		conf: setDefaults(viper.Sub("plugins." + name)),
	}
}

// Chan returns a strings channel with battery state monitor
func (b Battery) Chan() chan string {
	return b.out
}

// StartMonitor starts monitoring for battery changing events
func (b Battery) StartMonitor() {
	b.out <- b.formatMessage()
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(fmt.Errorf("cannot connect dbus session: %s", err.Error()))
	}

	arg := fmt.Sprintf(
		"type='signal',path='%s',interface='%s',member='%s',sender='%s'",
		fmt.Sprintf("/org/freedesktop/UPower/devices/battery_%s", b.conf.GetString("name")),
		"org.freedesktop.DBus.Properties",
		"PropertiesChanged",
		"org.freedesktop.UPower",
	)

	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		arg,
	)

	c := make(chan *dbus.Signal)
	conn.Signal(c)
	for {
		<-c
		b.out <- b.formatMessage()
	}
}

// formatMessage formats message for printing
func (b Battery) formatMessage() string {
	lvl, state := b.parseBatLevel()

	var pattern string
	if state != 2 {
		pattern = b.conf.GetString("ac")
	} else {
		switch {
		case lvl > 75:
			pattern = b.conf.GetString("high")
		case lvl > 35:
			pattern = b.conf.GetString("medium")
		case lvl > 12:
			pattern = b.conf.GetString("low")
		default:
			pattern = b.conf.GetString("empty")
		}
	}

	return utils.ReplaceVar(pattern, "lvl", strconv.Itoa(lvl))
}

// parseBatLevel connects to the system bus and get the State and Percentage
// properties from the UPower's BAT object. It returns the level in percents
// and integer status, which means:
//
//  0: Unknown
//  1: Charging
//  2: Discharging
//  3: Empty
//  4: Fully charged
//  5: Pending charge
//  6: Pending discharge
//
func (b Battery) parseBatLevel() (int, uint32) {
	conn, _ := dbus.SystemBus()
	pth := fmt.Sprintf("/org/freedesktop/UPower/devices/battery_%s", b.conf.GetString("name"))
	object := conn.Object(
		"org.freedesktop.UPower",
		dbus.ObjectPath(pth),
	)
	lvl, _ := object.GetProperty("org.freedesktop.UPower.Device.Percentage")
	state, _ := object.GetProperty("org.freedesktop.UPower.Device.State")
	return int(lvl.Value().(float64)), state.Value().(uint32)
}
