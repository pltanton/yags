package suspend

import (
	"fmt"
	"sync"

	"github.com/godbus/dbus"
)

var instance chan *dbus.Signal
var once = sync.Once{}

var subscribers = make([]chan *dbus.Signal, 0)

func monitorSuspend() {
	fmt.Println(subscribers)
	for {
		signalValue := <-instance
		for _, subscriber := range subscribers {
			subscriber <- signalValue
		}
	}
}

func startMonitorIfNeed() {
	once.Do(func() {
		if instance == nil {
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

			instance = make(chan *dbus.Signal, 1)
			conn.Signal(instance)

			go monitorSuspend()
		}
	})
}

func Subscribe(subscriber chan *dbus.Signal) {
	startMonitorIfNeed()
	subscribers = append(subscribers, subscriber)
}
