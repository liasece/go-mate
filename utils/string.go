package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func SnakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		// 判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			noUp := false
			// if i != num-1 && s[i+1] >= 'A' && s[i+1] <= 'Z' {
			// 	noUp = true
			// }
			if i > 0 && s[i-1] >= 'A' && s[i-1] <= 'Z' {
				noUp = true
			}
			if noUp && i < num-1 && s[i+1] >= 'a' && s[i+1] <= 'z' {
				noUp = false
			}
			if !noUp {
				data = append(data, '_')
			}
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	// ToLower把大写字母统一转小写
	return strings.ToLower(string(data))
}

func SnakeStringToBigHump(s string) string {
	ss := strings.Split(s, "_")
	for i, v := range ss {
		ss[i] = cases.Title(language.Und).String(v)
	}
	return strings.Join(ss, "")
}

func Title(s string) string {
	return cases.Title(language.Und).String(s)
}
