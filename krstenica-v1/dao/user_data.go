package dao

import (
	"regexp"
	"strings"
	"time"
)

var reUnderscore *regexp.Regexp

type HramDo struct {
	HramID    uint      `json:"hram_id"`
	HramName  string    `json:"hram_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func inList(elem string, list []string) bool {
	for _, el := range list {
		if el == elem {
			return true
		}
	}
	return false
}

// Underscore converts string from CamelCase to underscore case
func Underscore(in string) string {
	return strings.ToLower(reUnderscore.ReplaceAllStringFunc(in, convertToUnderscore))
}
func convertToUnderscore(s string) string {
	return string(s[0]) + "_" + string(s[1])
}
