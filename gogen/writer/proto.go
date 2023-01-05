package writer

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/liasece/go-mate/code"
	"github.com/liasece/gocoder"
	"github.com/liasece/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func StructToProto(protoFile string, indent string, ts ...gocoder.Type) error {
	if len(ts) == 0 {
		return nil
	}
	originFileContent := ""
	{
		// read from file
		content, err := os.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	parser := code.NewProtoBufCodeBlockParser()
	toCode := parser.Parse(originFileContent)
	for _, t := range ts {
		newContent := buildProtoContent(toCode.OriginString, t, indent)
		toCode.Merge(0, parser.Parse(newContent))
	}
	toContent := toCode.OriginString
	{
		// add end line
		if toContent != "" && toContent[len(toContent)-1] != '\n' {
			toContent += "\n"
		}
	}
	if toContent != originFileContent {
		// write to file
		err := os.WriteFile(protoFile, []byte(toContent), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

type ProtoInfo struct {
	Package string
}

func ReadProtoInfo(protoFile string) (*ProtoInfo, error) {
	// read from file
	content, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, err
	}

	var res ProtoInfo

	originContent := string(content)
	scanner := bufio.NewScanner(strings.NewReader(originContent))
	packageReg := regexp.MustCompile(`^\s*package\s+(\w+)\s*;?`)
	for scanner.Scan() {
		t := scanner.Text()
		if strings.HasPrefix(t, "package") {
			parts := packageReg.FindStringSubmatch(t)
			if len(parts) > 1 {
				res.Package = parts[1]
			}
		}
	}
	return &res, nil
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

func getProdFiledNumInOriginMsg(origin string, fieldNameRaw string) int {
	fieldReg := regexp.MustCompile(`.*?([0-9a-z_]+)\s*=\s*(\d+).*?`)
	parts := fieldReg.FindAllStringSubmatch(origin, -1)
	fieldName := strings.ReplaceAll(fieldNameRaw, "_", "")
	for _, fieldLine := range parts {
		if strings.ReplaceAll(fieldLine[1], "_", "") == fieldName {
			i, err := strconv.Atoi(fieldLine[2])
			if err != nil {
				log.Error("getProdFiledNumInOriginMsg Atoi error", log.ErrorField(err), log.Any("line", fieldLine))
			}
			return i
		}
	}
	return 0
}

func getMaxProdFiledNumInOriginMsg(origin string) int {
	fieldReg := regexp.MustCompile(`.*?([0-9a-z_]+)\s*=\s*(\d+).*?`)
	parts := fieldReg.FindAllStringSubmatch(origin, -1)
	res := 0
	for _, fieldLine := range parts {
		i, err := strconv.Atoi(fieldLine[2])
		if err != nil {
			log.Error("getProdFiledNumInOriginMsg Atoi error", log.ErrorField(err), log.Any("line", fieldLine))
		} else if i > res {
			res = i
		}
	}
	return res
}

func buildProtoContent(originContent string, t gocoder.Type, indent string) string {
	msgName := t.GetNamed()
	if msgName == "" {
		msgName = t.String()
	}
	matchOrigin := getProtoFromStr(originContent, msgName)
	fieldStr := make(map[int]string)
	maxOriginFieldNum := 0
	if matchOrigin != "" {
		maxOriginFieldNum = getMaxProdFiledNumInOriginMsg(matchOrigin)
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		typ := f.GetType().String()
		isBaseType := true
		tailPkg := f.GetType().AllSub()[0].PackageInReference()
		if strings.Count(typ, ".") == 0 && tailPkg != "" && tailPkg != "entity" {
			prefix := ""
			if strings.HasPrefix(typ, "*") {
				prefix = "*"
				typ = strings.ReplaceAll(typ, "*", "")
			}
			if strings.HasPrefix(typ, "[]") {
				prefix = "[]"
				typ = strings.ReplaceAll(typ, "[]", "")
			}
			typ = prefix + tailPkg + "." + typ
		}
		{
			ss := strings.Split(typ, ".")
			if len(ss) > 1 {
				for i, s := range ss {
					ss[i] = cases.Title(language.Und).String(s)
				}
				typ = strings.Join(ss, "")
				isBaseType = false
			}
		}
		isPtr := false
		if strings.Contains(typ, "*") {
			typ = strings.ReplaceAll(typ, "*", "")
			isPtr = true
		}
		isSlice := false
		if strings.Contains(typ, "[]") {
			typ = strings.ReplaceAll(typ, "[]", "")
			isSlice = true
		}
		switch typ {
		case "TimeTime":
			typ = "google.protobuf.Timestamp"
			isBaseType = false
		case "int":
			typ = "int64"
		case "float32":
			typ = "float"
		case "float64":
			typ = "double"
		}
		name := snakeString(f.GetName())
		opt := ""
		if isPtr {
			if isBaseType {
				opt = "optional "
			}
		}
		if isSlice {
			opt = "repeated "
		}
		if typ != "" {
			if originIndex := getProdFiledNumInOriginMsg(matchOrigin, name); originIndex > 0 {
				fieldStr[originIndex] = fmt.Sprintf("%s%s%s %s = %d;\n", indent, opt, typ, name, originIndex)
			} else {
				maxOriginFieldNum++
				fieldStr[maxOriginFieldNum] = fmt.Sprintf("%s%s%s %s = %d;\n", indent, opt, typ, name, maxOriginFieldNum)
			}
		} else {
			log.Debug("buildProtoContent skip type", log.Any("msgName", msgName), log.Any("name", name), log.Any("named", f.GetType().GetNamed()),
				log.Any("str", f.GetType().String()), log.Reflect("f", f))
		}
	}
	addFsStr := ""
	{
		is := make([]int, 0)
		for i := range fieldStr {
			is = append(is, i)
		}
		sort.Ints(is)
		for _, i := range is {
			addFsStr += fieldStr[i]
		}
	}
	return "message " + msgName + " {\n" + addFsStr + "}"
}
