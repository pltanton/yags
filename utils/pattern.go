package utils

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ReplaceVar returns a copy of the string s with replaced newValue substring,
// surrounded by the separators defined in config as varSeps, by newValue.
func ReplaceVar(s, varName, newValue string) string {
	varSeps := viper.GetString("varSeps")
	oldValue := fmt.Sprintf("%c%s%c", varSeps[0], varName, varSeps[1])
	return strings.Replace(s, oldValue, newValue, -1)
}
