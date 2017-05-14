package main

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("pause", 1000)
	v.SetDefault("connected", "{lvl}")
	v.SetDefault("disconnected", "N/A")
	v.SetDefault("interface", "wlan0")
	return v
}
