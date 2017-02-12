package timer

import "github.com/spf13/viper"

func setTimeDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("timeFormat", "Jan 2 15:04:05")
	v.SetDefault("pause", 1000)
	return v
}

func setWifiDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("pause", 1000)
	v.SetDefault("connected", "{lvl}")
	v.SetDefault("disconnected", "N/A")
	v.SetDefault("interface", "wlan0")
	return v
}
