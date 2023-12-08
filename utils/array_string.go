package utils

import (
	"fmt"
	"strings"
)

func arrayToString(value any) string {
	arrayStr := "{"

	if array, ok := value.([]any); ok {

		for _, val := range array {
			arrayStr += fmt.Sprintf("%v", val) + ","
		}

		return strings.TrimSuffix(arrayStr, ",") + "}"
	}

	return ""
}
