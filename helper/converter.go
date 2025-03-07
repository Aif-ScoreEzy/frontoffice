package helper

import (
	"fmt"
	"strconv"
)

func ConvertUintToString(arg uint) string {
	return strconv.Itoa(int(arg))
}

func InterfaceToUint(input interface{}) (uint, error) {
	if val, ok := input.(uint); ok {
		return val, nil
	}

	return 0, fmt.Errorf("cannot convert %T to uint", input)
}
