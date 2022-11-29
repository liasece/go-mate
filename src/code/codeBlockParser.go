package code

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type CodePairCount struct {
	Count   map[string]int // like "()"":1 "{}"":1 "[]"":1
	KeyWord []string       // like () {} [], in one string, first rune is left, second rune is right
}

func (c *CodePairCount) Add(line string) {
	for _, v := range c.KeyWord {
		head := ""
		tail := ""
		if len(v) == 2 {
			head = v[:1]
			tail = v[1:]
		} else {
			list := strings.Split(v, " ")
			if len(list) == 2 {
				head = list[0]
				tail = list[1]
			} else {
				panic("CodePairCount.Add: invalid keyword: " + v)
			}
		}
		if head == tail {
			if count := strings.Count(line, head); count > 0 {
				if c.Count[v] > 0 {
					c.Count[v] -= count % 2
				} else {
					c.Count[v] += count % 2
				}
			}
		} else {
			c.Count[v] += strings.Count(line, v[:1])
			c.Count[v] -= strings.Count(line, v[1:])
		}
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

type CodeBlockType struct {
	Name                   string
	MergeAble              bool
	RegStr                 *regexp.Regexp
	RegOriginIndex         int
	RegKeyIndex            int
	RegSubContentIndex     int
	ParentNames            []string
	SubsSeparator          string // like "," or ";"
	SubWarpChar            string // like "()" "{}" "[]"
	RegSubWarpContentIndex int
	KeyCaseIgnored         bool // ABc == abc
}

type CodeBlockParser struct {
	Types             []CodeBlockType
	PairKeys          []string
	LineCommentKey    string
	PendingLinePrefix string
}

func (c *CodeBlockParser) Parse(content string) *CodeBlock {
	res := &CodeBlock{
		Parent:          nil,
		OriginString:    content,
		SubOriginString: content,
		Type:            CodeBlockType{"", true, nil, 0, 0, 0, nil, "\n", "", 0, false},
	}
	res.SubList = c.protoBlocksFromString(res, res.SubOriginString)
	return res
}

func (c *CodeBlockParser) protoBlockFromString(parent *CodeBlock, content string) []*CodeBlock {
	for _, v := range c.Types {
		res := []*CodeBlock{}
		if v.RegStr == nil {
			continue
		}
		if len(v.ParentNames) > 0 {
			if parent == nil {
				continue
			}
			inParent := false
			for _, parentName := range v.ParentNames {
				if parentName == parent.Type.Name {
					inParent = true
					break
				}
			}
			if !inParent {
				continue
			}
		}
		// syntax
		contentReg := v.RegStr
		partsList := contentReg.FindAllStringSubmatch(content, -1)
		// if len(partsList) > 0 {
		// 	fmt.Println("protoBlockFromString:\n" + content + "\n" + fmt.Sprint(partsList))
		// }
		for _, parts := range partsList {
			item := &CodeBlock{
				Parent:       parent,
				OriginString: parts[v.RegOriginIndex],
			}
			item.Type = v
			if v.RegKeyIndex >= 0 {
				item.Key = parts[v.RegKeyIndex]
			}
			if v.RegSubContentIndex >= 0 {
				item.SubOriginString = parts[v.RegSubContentIndex]
				item.SubList = append(item.SubList, c.protoBlocksFromString(item, item.SubOriginString)...)
			}
			res = append(res, item)
		}
		if len(res) > 0 {
			return res
		}
	}
	{
		lineNoComment := strings.Split(content, c.LineCommentKey)[0]
		lineNoComment = strings.TrimSpace(lineNoComment)
		if lineNoComment != "" {
			fmt.Println("ProtoBlockFromString unknown:\n" + content)
		}
	}
	return nil
}

func (c *CodeBlockParser) protoBlocksFromString(parent *CodeBlock, content string) []*CodeBlock {
	// fmt.Println("ProtoBlocksFromString:\n" + content)
	scanner := bufio.NewScanner(strings.NewReader(content))
	res := make([]*CodeBlock, 0)
	pairCount := &CodePairCount{
		Count:   map[string]int{},
		KeyWord: c.PairKeys,
	}
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	currentBlockContent := ""
	for i, line := range lines {
		lineNoComment := strings.Split(line, c.LineCommentKey)[0]
		pairCount.Add(lineNoComment)
		currentBlockContent += line
		if i != len(lines)-1 || strings.HasSuffix(content, "\n") {
			currentBlockContent += "\n"
		}
		if pairCount.IsZero() {
			pending := false
			if c.PendingLinePrefix != "" {
				if len(lines) > i+1 {
					if strings.HasPrefix(strings.Trim(lines[i+1], " \t"), c.PendingLinePrefix) {
						pending = true
					}
				}
			}
			if !pending {
				res = append(res, c.protoBlockFromString(parent, currentBlockContent)...)
				currentBlockContent = ""
			}
		}
	}
	if currentBlockContent != "" {
		currentBlockContent = currentBlockContent[:len(currentBlockContent)-1]
	}
	if currentBlockContent != "" {
		res = append(res, c.protoBlockFromString(parent, currentBlockContent)...)
	}
	return res
}
