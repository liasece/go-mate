package utils

import "strings"

type TmplUtilsFunc struct {
}

func (TmplUtilsFunc) SplitN(origin string, sep string, n int) string {
	ss := strings.Split(origin, sep)
	if n < 0 || n >= len(ss) {
		return ""
	}
	return ss[n]
}

func (TmplUtilsFunc) ToLowerCamelCase(str string) string {
	if str == "" {
		return str
	}
	return strings.ToLower(str[:1]) + str[1:]
}

func (TmplUtilsFunc) ToCamelCase(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}
