package battery

import "github.com/spf13/viper"

func setDefaults(v *viper.Viper) *viper.Viper {
	v.SetDefault("name", "BAT0")
	v.SetDefault("high", "{lvl}")
	v.SetDefault("medium", "{lvl}")
	v.SetDefault("low", "{lvl}")
	v.SetDefault("empty", "{lvl}")
	v.SetDefault("ac", "{lvl}")
	return v
}
