package code

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"
)

type CodePairCount struct {
	Count      map[string]int // like "()":1 "{}":1 "[]":1
	KeyWord    []string       // like () {} [], in one string, first rune is left, second rune is right
	OriginText []string
}

func (c *CodePairCount) Add(line string) {
	inOrigin := ""
	for k, v := range c.Count {
		if v > 0 {
			for _, originWarp := range c.OriginText {
				if k == originWarp {
					inOrigin = originWarp
					break
				}
			}
		}
	}
	for _, key := range c.KeyWord {
		if inOrigin != "" {
			if key != inOrigin {
				continue
			}
		}
		head := ""
		tail := ""
		if len(key) == 2 {
			head = key[:1]
			tail = key[1:]
		} else {
			list := strings.Split(key, " ")
			if len(list) == 2 {
				head = list[0]
				tail = list[1]
			} else {
				panic("CodePairCount.Add: invalid keyword: " + key)
			}
		}
		if head == tail {
			if count := strings.Count(line, head); count > 0 {
				if c.Count[key] > 0 {
					c.Count[key] -= count % 2
				} else {
					c.Count[key] += count % 2
				}
			}
		} else {
			c.Count[key] += strings.Count(line, head)
			c.Count[key] -= strings.Count(line, tail)
		}

		// no repeated count
		line = strings.ReplaceAll(line, head, "")
		line = strings.ReplaceAll(line, tail, "")
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
	SubsSeparator          string // the sub body separator character, like "," or ";"
	SubWarpChar            string // the sub body warp character, like "()" "{}" "[]" "((|))" `"""|"""`
	RegSubWarpContentIndex int
	KeyCaseIgnored         bool     // ABc == abc
	SubTailChar            []string // the sub body tail character, like "}" ";" "))" ","
}

type CodeBlockParser struct {
	Types             []CodeBlockType
	PairKeys          []string
	OriginText        []string
	LineCommentKey    string
	PendingLinePrefix string
}

func (c *CodeBlockParser) Parse(content string) *CodeBlock {
	res := &CodeBlock{
		Parent:          nil,
		OriginString:    content,
		SubOriginString: content,
		Type:            CodeBlockType{"", true, nil, 0, 0, 0, nil, "\n", "", 0, false, nil},
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
		Count:      map[string]int{},
		KeyWord:    c.PairKeys,
		OriginText: c.OriginText,
	}
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	currentBlockContent := ""
	for i, line := range lines {
		lineNoComment := strings.Split(line, c.LineCommentKey)[0]
		pairCount.Add(lineNoComment)
		{
			// append code line
			hasContent := true
			if currentBlockContent == "" && strings.TrimSpace(lineNoComment) == "" {
				hasContent = false
			}
			if hasContent {
				currentBlockContent += line
				if i != len(lines)-1 || strings.HasSuffix(content, "\n") {
					currentBlockContent += "\n"
				}
			}
		}
		pending := false
		{
			// check need pending line
			if c.PendingLinePrefix != "" {
				if len(lines) > i+1 && strings.HasPrefix(strings.Trim(lines[i+1], " \t"), c.PendingLinePrefix) {
					pending = true
				}
			}
			if !pending && len(parent.Type.SubTailChar) > 0 {
				// must end with tail char
				find := false
				for _, tail := range parent.Type.SubTailChar {
					if strings.HasSuffix(strings.Trim(line, " \t"), tail) {
						find = true
						break
					}
				}
				if !find {
					pending = true
				}
			}
		}
		if !pending {
			if pairCount.IsZero() {
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
