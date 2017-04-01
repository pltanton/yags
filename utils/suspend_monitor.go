package utils

import (
	"fmt"

	"github.com/godbus/dbus"
)

// GetResumeDbusConn returns DBUS connection to resume from suspend signal
func GetResumeDbusConn() *dbus.Conn {
	conn, err := dbus.SystemBus()
	if err != nil {
		panic(fmt.Errorf("cannot connect dbus session: %s", err.Error()))
	}

	arg := fmt.Sprintf(
		"type='signal',interface='%s',member='%s',sender='%s'",
		"org.freedesktop.login1.Manager",
		"PrepareForSleep",
		"org.freedesktop.login1",
	)

	conn.BusObject().Call(
		"org.freedesktop.DBus.AddMatch",
		0,
		arg,
	)

	return conn
}
