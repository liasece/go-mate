package utils

import (
	"reflect"
	"strings"
)

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func TagMatch(tagReg string, tag string) bool {
	// filterSS := strings.Split(tagReg, ":")
	// tagName := strings.TrimSpace(filterSS[0])
	// filterValue := ""
	// if len(filterSS) > 1 {
	// 	filterValue = strings.TrimSpace(filterSS[1])
	// }
	// tagValue := reflect.StructTag(tag).Get(tagName)
	// if tagValue == "" {
	// 	return false
	// }
	// if filterValue == "" {
	// 	log.L(nil).Error("TagMatch: filterValue is empty", log.String("tagReg", tagReg), log.String("tag", tag))
	// 	return true
	// }
	// values := strings.Split(tagValue, ",")
	// for _, v := range values {
	// 	if strings.TrimSpace(v) == filterValue {
	// 		log.L(nil).Error("TagMatch: got", log.String("tagReg", tagReg), log.String("tag", tag), log.String("v", v))
	// 		return true
	// 	}
	// }
	return GetTagMatch(tagReg, tag) != ""
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`, return `url`
// tagReg like: `gomate:url` match: `gomate:"foo,url:testURL"`, return `testURL`
func GetTagMatch(tagReg string, tag string) string {
	filterSS := strings.Split(tagReg, ":")
	tagName := strings.TrimSpace(filterSS[0])
	filterValue := ""
	if len(filterSS) > 1 {
		filterValue = strings.TrimSpace(filterSS[1])
	}
	tagValue := reflect.StructTag(tag).Get(tagName)
	if tagValue == "" {
		return ""
	}
	if filterValue == "" {
		return tagValue
	}
	values := strings.Split(tagValue, ",")
	for _, v := range values {
		subValues := strings.Split(v, ":")
		if strings.TrimSpace(subValues[0]) == filterValue {
			if len(subValues) > 1 {
				return strings.TrimSpace(subValues[1])
			} else {
				return subValues[0]
			}
		}
	}
	return ""
}
