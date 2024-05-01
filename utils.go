package main

import "strings"

func ParseStringToBool(s string) bool {
	s = strings.ToLower(s)
	switch s {
	case "false", "0", "":
		return false
	default:
		return true
	}
}
