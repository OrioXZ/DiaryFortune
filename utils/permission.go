package utils

import "strings"

func IsAdmin(username string) bool {
	return strings.TrimSpace(strings.ToLower(username)) == "admin"
}
