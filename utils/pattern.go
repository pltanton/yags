package utils

import (
	"fmt"
	"regexp"
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

// GetVars returns the variables containing in the stirng s
func GetVars(s string) (result []string) {
	seps := viper.GetString("varSeps")
	regex := regexp.MustCompile(fmt.Sprintf(`%c([[:word:]]+)%c`, seps[0], seps[1]))
	allSubmatches := regex.FindAllStringSubmatch(s, -1)
	for _, submatch := range allSubmatches {
		result = append(result, submatch[1])
	}
	return
}

// Contains return true if key is inside slice
func Contains(key string, arr *[]string) bool {
	for _, val := range *arr {
		if key == val {
			return true
		}
	}
	return false
}
