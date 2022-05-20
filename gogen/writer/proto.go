package writer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/liasece/gocoder"
)

func StructToProto(protoFile string, t gocoder.Type, indent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := ioutil.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := buildProtoContent(originFileContent, t, indent)
	{
		// write to file
		err := ioutil.WriteFile(protoFile, []byte(toContent), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func getProtoFromStr(originContent string, typ string) string {
	scanner := bufio.NewScanner(strings.NewReader(originContent))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	out := ""
	in := false
	for _, l := range lines {
		if !in {
			nameReg := regexp.MustCompile(`^\s*message\s+` + typ + `\s*{`)
			parts := nameReg.FindStringSubmatch(l)
			if len(parts) != 0 {
				in = true
			}
		}
		if in {
			out += l

			nameReg := regexp.MustCompile(`^\s*}`)
			parts := nameReg.FindStringSubmatch(l)
			if len(parts) != 0 {
				break
			}
			out += "\n"
		}
	}
	return out
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
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
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

func buildProtoContent(originContent string, t gocoder.Type, indent string) string {
	msgName := t.GetNamed()
	if msgName == "" {
		msgName = t.String()
	}
	matchOrigin := getProtoFromStr(originContent, msgName)
	addFsStr := ""
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		// log.Error("in buildProtoContent filed", log.Any("i", i), log.Any("f", f))
		typ := f.GetType().GetNamed()
		isBaseType := true
		ss := strings.Split(typ, ".")
		for i, s := range ss {
			ss[i] = strings.Title(s)
		}
		typ = strings.Join(ss, "")
		if typ == "" {
			typ = f.GetType().String()
		} else {
			isBaseType = false
		}
		isPtr := false
		if strings.HasPrefix(typ, "*") {
			typ = typ[1:]
			isPtr = true
		}
		isSlice := false
		if strings.HasPrefix(typ, "[]") {
			typ = typ[2:]
			isSlice = true
		}
		switch typ {
		case "time.Time":
			typ = "google.protobuf.Timestamp"
			isBaseType = false
		case "int":
			typ = "int64"
		}
		name := snakeString(f.GetName())
		opt := ""
		if isBaseType {
			if isPtr {
				opt = "optional "
			}
			if isSlice {
				opt = "repeated "
			}
		}
		addFsStr += fmt.Sprintf("%s%s%s %s = %d;\n", indent, opt, typ, name, i+1)
	}
	toStr := "message " + msgName + " {\n" + addFsStr + "}"
	if matchOrigin != "" {
		return strings.Replace(originContent, matchOrigin, toStr, 1)
	}
	return originContent + "\n" + toStr + "\n"
}
