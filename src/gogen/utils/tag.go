package utils

import (
	"reflect"
	"strings"
)

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func TagMatch(tagReg string, tag string) bool {
	filterSS := strings.Split(tagReg, ":")
	tagName := strings.TrimSpace(filterSS[0])
	filterValue := ""
	if len(filterSS) > 1 {
		filterValue = strings.TrimSpace(filterSS[1])
	}
	tagValue := reflect.StructTag(tag).Get(tagName)
	if tagValue == "" {
		return false
	}
	if filterValue == "" {
		return true
	}
	values := strings.Split(tagValue, ",")
	for _, v := range values {
		if strings.TrimSpace(v) == filterValue {
			return true
		}
	}
	return false
}
