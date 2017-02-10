package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/core"
)

func main() {
	if len(os.Args) != 2 {
		panic(fmt.Errorf("You should specify the config file as only one argument"))
	}
	viper.SetConfigFile(os.Args[1])
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	core.Init()
}
