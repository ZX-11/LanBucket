package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

const (
	port = ":18800"
)

var settings = map[string]interface{}{
	"EnableUpload": false,
	"EnableIPv6":   false,
}

func autoSet(cmd string) error {
	args := strings.SplitN(cmd, " ", 2)
	if len(args) != 2 {
		return fmt.Errorf("Error: Arguement not found.")
	}
	return set(args[0], parse(args[1]))
}

func set(item string, value interface{}) error {
	oldValue, ok := settings[item]
	if !ok {
		return fmt.Errorf(`Error: Setting "%v" not found.`, item)
	}
	if reflect.TypeOf(value).Kind() != reflect.TypeOf(oldValue).Kind() {
		return fmt.Errorf("Error: Type mismatched.")
	}
	settings[item] = value
	return nil
}

func parse(arg string) interface{} {
	isNumber := true
	for _, v := range arg {
		if !unicode.IsNumber(v) {
			isNumber = false
			break
		}
	}
	switch {
	case isNumber:
		val, _ := strconv.ParseUint(arg, 10, 64)
		return val
	case arg == "true":
		return true
	case arg == "false":
		return false
	default:
		return strings.Trim(arg, `"`)
	}
}
