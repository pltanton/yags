package core

import "github.com/spf13/viper"

func setDefaults() {
	viper.SetDefault("varSeps", "{}")
}
