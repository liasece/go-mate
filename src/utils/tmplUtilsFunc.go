package utils

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TmplUtilsFunc struct {
}

func (TmplUtilsFunc) SplitN(origin string, sep string, n int) string {
	ss := strings.Split(origin, sep)
	if n < 0 {
		if -1*n > len(ss) {
			return ""
		}
		return ss[len(ss)+n]
	} else {
		if n >= len(ss) {
			return ""
		}
		return ss[n]
	}
}

func (t TmplUtilsFunc) ToLowerCamelCase(str string) string {
	if str == "" {
		return str
	}
	words := t.SplitCamelCase(str)
	return strings.ToLower(words[0]) + strings.Join(words[1:], "")
}

func (TmplUtilsFunc) ToCamelCase(str string) string {
	if str == "" {
		return str
	}
	return strings.ToUpper(str[:1]) + str[1:]
}

func (TmplUtilsFunc) HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

func (TmplUtilsFunc) HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// english word to plural
func (t TmplUtilsFunc) Plural(word string) string {
	if word == "" {
		return word
	}
	if strings.HasSuffix(word, "y") {
		return word[:len(word)-1] + "ies"
	}
	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") || strings.HasSuffix(word, "o") || strings.HasSuffix(word, "ch") || strings.HasSuffix(word, "sh") {
		return word + "es"
	}
	return word + "s"
}

// Split splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
// For more info please check: http://en.wikipedia.org/wiki/CamelCase
//
// Examples
//
//   "" =>                     []
//   "lowercase" =>            ["lowercase"]
//   "Class" =>                ["Class"]
//   "MyClass" =>              ["My", "Class"]
//   "MyC" =>                  ["My", "C"]
//   "HTML" =>                 ["HTML"]
//   "PDFLoader" =>            ["PDF", "Loader"]
//   "AString" =>              ["A", "String"]
//   "SimpleXMLParser" =>      ["Simple", "XML", "Parser"]
//   "vimRPCPlugin" =>         ["vim", "RPC", "Plugin"]
//   "GL11Version" =>          ["GL", "11", "Version"]
//   "99Bottles" =>            ["99", "Bottles"]
//   "May5" =>                 ["May", "5"]
//   "BFG9000" =>              ["BFG", "9000"]
//   "BöseÜberraschung" =>     ["Böse", "Überraschung"]
//   "Two  spaces" =>          ["Two", "  ", "spaces"]
//   "BadUTF8\xe2\xe2\xa1" =>  ["BadUTF8\xe2\xe2\xa1"]
//
// Splitting rules
//
//  1) If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2) Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3) Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4) Iterate through array of split strings, and if a given string
//     is upper case:
//       if subsequent string is lower case:
//         move last character of upper case string to beginning of
//         lower case string
func (t TmplUtilsFunc) SplitCamelCase(src string) (entries []string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		return []string{src}
	}
	entries = []string{}
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	return
}

func (t TmplUtilsFunc) GraphqlStyle(fieldName string, typeName string) string {
	return t.GraphqlStyleW(fieldName, typeName, "")
}

func (t TmplUtilsFunc) GraphqlStyleW(fieldName string, typeName string, writeType string) string {
	nameWords := t.SplitCamelCase(fieldName)
	isID := false
	for _, word := range nameWords {
		if strings.ToLower(word) == "id" {
			isID = true
			break
		}
	}
	isSlice := false
	if strings.HasPrefix(typeName, "[]") {
		isSlice = true
		typeName = typeName[2:]
	}
	isPtr := false
	if strings.HasPrefix(typeName, "*") {
		isPtr = true
		typeName = typeName[1:]
	}
	res := ""
	switch typeName {
	case "string":
		if isID {
			res = "ID"
		} else {
			res = "String"
		}
	case "int", "int32", "int64", "uint", "uint32", "uint64":
		res = "Int"
	case "float32", "float64":
		res = "Float"
	case "time.Time":
		res = "Time"
	case "bool":
		res = "Boolean"
	default:
		if writeType != "" && typeName == writeType {
			contentReg := regexp.MustCompile(`^[A-Z][\w]*$`)
			if contentReg.MatchString(typeName) {
				isPtr = false
			}
			res = typeName
		} else {
			return ""
		}
	}
	if !isPtr {
		res = res + "!"
	}
	if isSlice {
		res = "[" + res + "]"
	}
	return res
}

func (t TmplUtilsFunc) ReplaceWord(str string, old string, new string) string {
	nameWords := t.SplitCamelCase(str)
	for i, word := range nameWords {
		if word == old {
			nameWords[i] = new
		}
	}
	return strings.Join(nameWords, "")
}

func (t TmplUtilsFunc) ReplaceWord2(str string, old string, new string, old2 string, new2 string) string {
	nameWords := t.SplitCamelCase(str)
	for i, word := range nameWords {
		if word == old {
			nameWords[i] = new
		}
		if word == old2 {
			nameWords[i] = new2
		}
	}
	return strings.Join(nameWords, "")
}
