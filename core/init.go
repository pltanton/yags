package core

import (
	"fmt"
	"reflect"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/plugins"
	"github.com/pltanton/yags/plugins/battery"
	"github.com/pltanton/yags/plugins/kbdd"
	"github.com/pltanton/yags/plugins/stdin"
	"github.com/pltanton/yags/plugins/timer"
	"github.com/pltanton/yags/plugins/volume"
	"github.com/pltanton/yags/utils"
)

var pluginsNames []string
var pluginsInstances []plugins.Plugin
var cases []reflect.SelectCase
var values map[string]string

func initPlugins() {
	values = make(map[string]string)
	pluginsNames = utils.GetVars(viper.GetString("format"))
	pluginsInstances = make([]plugins.Plugin, len(pluginsNames))

	for i, name := range pluginsNames {
		var plugin plugins.Plugin
		if name == "stdin" {
			plugin = stdin.NewStdin()
		} else {
			typ := viper.GetString("plugins." + name + ".type")
			switch typ {
			case "volume":
				plugin = volume.NewVolume(name)
			case "battery":
				plugin = battery.NewBattery(name)
			case "timer":
				plugin = timer.NewTimerCMD(name)
			case "time":
				plugin = timer.NewTime(name)
			case "wifi":
				plugin = timer.NewWifi(name)
			case "kbdd":
				plugin = kbdd.NewKBDD(name)
			default:
				continue
			}
		}

		pluginsInstances[i] = plugin

		go plugin.StartMonitor()
	}
}

func initCases() {
	cases = make([]reflect.SelectCase, len(pluginsInstances))
	for i, pluginInstance := range pluginsInstances {
		cases[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(pluginInstance.Chan()),
		}
	}
}

func listen() {
	for {
		chosen, value, _ := reflect.Select(cases)
		values[pluginsNames[chosen]] = value.String()
		fmt.Println(formatOutput())
	}
}

func Init() {
	setDefaults()
	initPlugins()
	initCases()
	listen()
}
