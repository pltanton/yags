package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/core"
)

func main() {
	viper.SetConfigName(os.Args[1])
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	core.Init()
}
