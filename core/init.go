package core

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/pltanton/yags/plugins"
	"github.com/pltanton/yags/utils"
	"github.com/spf13/viper"
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
		path := viper.GetString("plugins." + name + ".path")
		path, _ = filepath.Abs(os.ExpandEnv(path))

		p, err := plugin.Open(path)
		if err != nil {
			panic(err)
		}

		fmt.Println("here")

		pluginNewSym, err := p.Lookup("New")
		if err != nil {
			panic(err)
		}

		pluginNew := pluginNewSym.(func(conf *viper.Viper) plugins.Plugin)

		pluginConfig := viper.Sub("plugins." + name)
		pluginInstance := pluginNew(pluginConfig)

		pluginsInstances[i] = pluginInstance

		go pluginInstance.StartMonitor()
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
