package utils

import (
	"fmt"
	"strings"
)

func Interpolate(template string, params map[string]interface{}) string {
	var result string = template
	for name, param := range params {
		placeholder := fmt.Sprintf("{%s}", name)
		stringValue := fmt.Sprintf("%v", param)
		result = strings.Replace(result, placeholder, stringValue, -1)
	}
	return result
}
