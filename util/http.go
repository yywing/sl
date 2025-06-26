package util

import "strings"

func HttpSchemeToPort(scheme string) int {
	if strings.ToLower(scheme) == "http" {
		return 80
	} else if strings.ToLower(scheme) == "https" {
		return 443
	}
	return 0
}
