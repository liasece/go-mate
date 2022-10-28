package writer

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
)

type CodePairCount struct {
	Count   map[string]int // like "()"":1 "{}"":1 "[]"":1
	KeyWord []string       // like () {} [], in one string, first rune is left, second rune is right
}

func (c *CodePairCount) Add(line string) {
	for _, v := range c.KeyWord {
		c.Count[v] += strings.Count(line, v[:1])
		c.Count[v] -= strings.Count(line, v[1:])
	}
}

func (c *CodePairCount) IsZero() bool {
	for _, v := range c.KeyWord {
		if c.Count[v] != 0 {
			return false
		}
	}
	return true
}

type ProtoBlockType string

const (
	ProtoBlockTypeNone         ProtoBlockType = ""
	ProtoBlockTypeImport       ProtoBlockType = "import"
	ProtoBlockTypePackage      ProtoBlockType = "package"
	ProtoBlockTypeSyntax       ProtoBlockType = "syntax"
	ProtoBlockTypeService      ProtoBlockType = "service"
	ProtoBlockTypeRPC          ProtoBlockType = "rpc"
	ProtoBlockTypeMessage      ProtoBlockType = "message"
	ProtoBlockTypeMessageField ProtoBlockType = "message field"
	ProtoBlockTypeOption       ProtoBlockType = "options"
	ProtoBlockTypeOptionKey    ProtoBlockType = "option item"
	ProtoBlockTypeEnum         ProtoBlockType = "enum"
	ProtoBlockTypeEnumKey      ProtoBlockType = "enum item"
	ProtoBlockTypeReserved     ProtoBlockType = "reserved"
)

type ProtoBlock struct {
	Parent          *ProtoBlock `json:"-"`
	Key             string
	Type            ProtoBlockType
	OriginString    string
	SubOriginString string
	SubList         []*ProtoBlock
}

func (b *ProtoBlock) Find(typ ProtoBlockType, key string) *ProtoBlock {
	if b.Type == typ && b.Key == key {
		return b
	}
	for _, v := range b.SubList {
		if v.Type == typ && v.Key == key {
			return v
		}
	}
	return nil
}

func (b *ProtoBlock) addSub(income *ProtoBlock) {
	if len(b.SubList) == 0 {
		panic("addSub to empty block: " + income.Key)
	}
	b.SubList = append(b.SubList, income)
	newSubOriginString := fmt.Sprintf("%s%s", b.SubOriginString, income.OriginString)
	switch income.Type {
	case ProtoBlockTypeOptionKey:
		if !strings.HasSuffix(strings.TrimSpace(b.SubOriginString), ",") {
			newSubOriginString = fmt.Sprintf("%s, %s", b.SubOriginString, income.OriginString)
		}
	}
	for b != nil {
		myOldSubOriginString := b.SubOriginString
		myOldOriginString := b.OriginString

		// update sub origin
		b.SubOriginString = newSubOriginString

		// update origin
		b.OriginString = strings.Replace(b.OriginString, myOldSubOriginString, b.SubOriginString, 1)
		// fmt.Println(b.Key + " new OriginString: " + b.OriginString + " (replaceFrom: " + myOldSubOriginString + ")" + " (replaceTo: " + b.SubOriginString + ")")

		if b.Parent != nil {
			newSubOriginString = strings.Replace(b.Parent.SubOriginString, myOldOriginString, b.OriginString, 1)
		}
		b = b.Parent
	}
}

func (b *ProtoBlock) Merge(income *ProtoBlock) *ProtoBlock {
	exists := b.Find(income.Type, income.Key)
	if exists == nil {
		// fmt.Println("Merge not exists:" + income.Key)
		// append to current block
		b.addSub(income)
	} else {
		// fmt.Println("Merge exists:" + income.Key)
		switch exists.Type {
		case ProtoBlockTypeService, ProtoBlockTypeMessage, ProtoBlockTypeOption, ProtoBlockTypeNone, ProtoBlockTypeMessageField:
			for _, v := range income.SubList {
				exists.Merge(v)
			}
		}
	}
	return b
}

func ProtoBlockFromString(content string) *ProtoBlock {
	res := &ProtoBlock{
		Parent:          nil,
		OriginString:    content,
		SubOriginString: content,
	}
	res.SubList = ProtoBlocksFromString(res, res.SubOriginString)
	return res
}

func protoBlockFromString(p *ProtoBlock, content string) *ProtoBlock {
	// fmt.Println("ProtoBlockFromString:\n" + content)
	res := &ProtoBlock{
		Parent:       p,
		OriginString: content,
	}
	getOptionItemsStr := func(content string) [][]string {
		contentReg := regexp.MustCompile(`(?s)\s*(\w+)\s*[:=]\s*(\S+)\s*`)
		return contentReg.FindAllStringSubmatch(content, -1)
	}
	getEnumItemsStr := func(content string) [][]string {
		contentReg := regexp.MustCompile(`(?s)\s*(\w+)\s*[:=]\s*(\S+)\s*`)
		return contentReg.FindAllStringSubmatch(content, -1)
	}
	{
		// service
		contentReg := regexp.MustCompile(`(?s)^\s*service\s+(\w+)\s*\{(.*?)\}\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeService
			res.Key = parts[1]
			res.SubOriginString = parts[2]
			res.SubList = append(res.SubList, ProtoBlocksFromString(res, res.SubOriginString)...)
			return res
		}
	}
	{
		// rpc
		contentReg := regexp.MustCompile(`(?s)^\s*rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{(.*)\})?\s*;?(\s*//.*)?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeRPC
			res.Key = parts[1]
			res.SubOriginString = parts[3]
			res.SubList = append(res.SubList, ProtoBlocksFromString(res, res.SubOriginString)...)
			return res
		}
	}
	{
		// message
		contentReg := regexp.MustCompile(`(?s)^\s*message\s+(\w+)\s*\{(.*)\}\s*;?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeMessage
			res.Key = parts[1]
			res.SubOriginString = parts[2]
			res.SubList = append(res.SubList, ProtoBlocksFromString(res, res.SubOriginString)...)
			return res
		}
	}
	if p != nil && p.Type == ProtoBlockTypeMessage {
		// message field
		contentReg := regexp.MustCompile(`(?s)^\s*(optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?(\s*//.*)?\s*?$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeMessageField
			res.Key = parts[4]
			{
				// check option key
				res.SubOriginString = parts[7]
				partsList := getOptionItemsStr(res.SubOriginString)
				if len(partsList) > 0 {
					for _, parts := range partsList {
						block := protoBlockFromString(res, parts[0])
						if block != nil {
							res.SubList = append(res.SubList, block)
						}
					}
				}
			}
			return res
		}
	}
	{
		// origin option
		contentReg := regexp.MustCompile(`(?s)^\s*option\s+(.+?)\s*=\s*\{?(.*)\}?\s*;?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeOption
			res.Key = parts[1]
			{
				// check option key
				res.SubOriginString = parts[2]
				partsList := getOptionItemsStr(res.SubOriginString)
				if len(partsList) > 0 {
					for _, parts := range partsList {
						block := protoBlockFromString(res, parts[0])
						if block != nil {
							res.SubList = append(res.SubList, block)
						}
					}
				}
			}
			return res
		}
	}
	if p != nil && (p.Type == ProtoBlockTypeOption || p.Type == ProtoBlockTypeMessageField) {
		// option key
		partsList := getOptionItemsStr(content)
		for _, parts := range partsList {
			if len(parts) > 0 {
				res.Type = ProtoBlockTypeOptionKey
				res.Key = parts[1]
				return res
			}
		}
	}
	{
		// enum
		contentReg := regexp.MustCompile(`(?s)^\s*enum\s*(.+?)\s*\{(.*)\}\s*;?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeEnum
			res.Key = parts[1]
			{
				// check option key
				res.SubOriginString = parts[2]
				partsList := getEnumItemsStr(res.SubOriginString)
				if len(partsList) > 0 {
					for _, parts := range partsList {
						block := protoBlockFromString(res, parts[0])
						if block != nil {
							res.SubList = append(res.SubList, block)
						}
					}
				} else {
					res.SubList = append(res.SubList, ProtoBlocksFromString(res, res.SubOriginString)...)
				}
			}
			return res
		}
	}

	if p != nil && (p.Type == ProtoBlockTypeEnum || p.Type == ProtoBlockTypeMessageField) {
		// enum item
		partsList := getEnumItemsStr(content)
		for _, parts := range partsList {
			if len(parts) > 0 {
				res.Type = ProtoBlockTypeEnumKey
				res.Key = parts[1]
				return res
			}
		}
	}
	{
		// reserved
		contentReg := regexp.MustCompile(`(?s)^\s*reserved\s*(.+?);?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeReserved
			res.Key = parts[1]
			return res
		}
	}
	{
		// import
		contentReg := regexp.MustCompile(`(?s)^\s*import\s*(.+?);?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeImport
			res.Key = parts[1]
			return res
		}
	}
	{
		// package
		contentReg := regexp.MustCompile(`(?s)^\s*package\s*(.+?);?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypePackage
			res.Key = parts[1]
			return res
		}
	}
	{
		// syntax
		contentReg := regexp.MustCompile(`(?s)^\s*syntax\s*=\s*(.+?);?\s*$`)
		parts := contentReg.FindStringSubmatch(content)
		if len(parts) > 0 {
			res.Type = ProtoBlockTypeSyntax
			res.Key = parts[1]
			return res
		}
	}
	{
		lineNoComment := strings.Split(content, "//")[0]
		lineNoComment = strings.TrimSpace(lineNoComment)
		if lineNoComment != "" {
			fmt.Println("ProtoBlockFromString unknown:\n" + content)
		}
	}
	return nil
}

func ProtoBlocksFromString(p *ProtoBlock, content string) []*ProtoBlock {
	// fmt.Println("ProtoBlocksFromString:\n" + content)
	scanner := bufio.NewScanner(strings.NewReader(content))
	res := make([]*ProtoBlock, 0)
	pairCount := &CodePairCount{
		Count:   map[string]int{},
		KeyWord: []string{"{}", "[]", "()"},
	}
	currentBlockContent := ""
	for scanner.Scan() {
		line := scanner.Text()
		lineNoComment := strings.Split(line, "//")[0]
		pairCount.Add(lineNoComment)
		if pairCount.IsZero() {
			currentBlockContent += line + "\n"
			block := protoBlockFromString(p, currentBlockContent)
			if block != nil {
				res = append(res, block)
			}
			currentBlockContent = ""
		} else {
			currentBlockContent += line + "\n"
		}
	}
	if currentBlockContent != "" {
		block := protoBlockFromString(p, currentBlockContent)
		if block != nil {
			res = append(res, block)
		}
	}
	return res
}

func MergeProtoFromFile(protoFile string, newContent string) error {
	return mergeProtoFromFile(protoFile, newContent)
}

func mergeProtoFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := ioutil.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeProto(originFileContent, newContent)
	if toContent != originFileContent {
		// write to file
		err := ioutil.WriteFile(protoFile, []byte(toContent), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeProto(originContent string, newContent string) string {
	res := ProtoBlockFromString(originContent).Merge(ProtoBlockFromString(newContent))
	return res.OriginString
}
