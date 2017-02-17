package volume

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("high", "{vol}")
	v.SetDefault("medium", "{vol}")
	v.SetDefault("low", "{vol}")
	v.SetDefault("muted", "{vol}")
	v.SetDefault("sink", "/org/pulseaudio/core1/sink0")
	return v
}
