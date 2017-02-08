package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	"github.com/pltanton/yags/plugins/volume"
)

func main() {
	viper.SetConfigName(os.Args[1])
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	volume.Monitor()
}
