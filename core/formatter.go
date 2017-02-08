package core

import (
	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

func formatOutput() string {
	formatString := viper.GetString("format")
	for _, name := range pluginsNames {
		formatString = utils.ReplaceVar(
			formatString,
			name,
			values[name],
		)
	}
	return formatString
}
