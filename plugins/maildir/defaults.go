package main

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("empty", "{lvl}")
	v.SetDefault("filled", "{lvl}")
	return v
}
