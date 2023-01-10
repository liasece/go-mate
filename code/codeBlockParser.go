package code

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/rs/xid"
)

type BlockParser struct {
	Types              []BlockType
	PairKeys           []string
	OriginText         []string
	HeadViscousPairKey []string // this PairKey will viscous the next block
	LineCommentKey     string
	PendingLinePrefix  string
}

func (c *BlockParser) Parse(content string) *Block {
	res := &Block{
		ID:               xid.New().String(),
		BlockParser:      c,
		Key:              "",
		SubList:          nil,
		Parent:           nil,
		OriginString:     content,
		SubOriginString:  []string{content},
		RegOriginStrings: nil,
		RegOriginIndexes: nil,
		SubOriginIndex:   nil,
		Type:             BlockType{"", nil, 0, 0, nil, nil, []*MergeConfig{{Append: true, ReplaceBlockType: nil}}, nil, "\n", "", 0, false, nil},
	}
	res.SubList = [][]*Block{c.protoBlocksFromString(res, res.SubOriginString[0], nil)}
	return res
}

func (c *BlockParser) protoBlockFromString(parent *Block, content string, mustTypeName []string) []*Block {
	allowTypeNames := make([]string, 0)
	for _, v := range c.Types {
		res := []*Block{}
		if v.RegStr == nil {
			continue
		}
		if mustTypeName != nil {
			find := false
			for _, name := range mustTypeName {
				if v.Name == name {
					find = true
					break
				}
			}
			if !find {
				continue
			}
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

		allowTypeNames = append(allowTypeNames, v.Name)

		// syntax
		contentReg := v.RegStr
		partsList := contentReg.FindAllStringSubmatch(content, -1)
		partsIndexList := contentReg.FindAllStringSubmatchIndex(content, -1)
		// log.Warn("partsIndexList", log.Any("partsIndexList", partsIndexList), log.Any("partsList", partsList))
		for partIndex, parts := range partsList {
			item := &Block{
				ID:               xid.New().String(),
				BlockParser:      c,
				Parent:           parent,
				OriginString:     parts[0],
				Key:              "",
				Type:             v,
				RegOriginStrings: make([]string, len(parts)),
				RegOriginIndexes: make([][]int, len(parts)),
				SubOriginString:  make([]string, len(v.RegSubContentIndex)),
				SubOriginIndex:   make([][]int, len(v.RegSubContentIndex)),
				SubList:          make([][]*Block, len(v.RegSubContentIndex)),
			}
			for i, v := range parts {
				item.RegOriginStrings[i] = v
				index1, index2 := partsIndexList[partIndex][i*2], partsIndexList[partIndex][i*2+1]
				if index1 >= 0 || index2 >= 0 {
					item.RegOriginIndexes[i] = []int{index1, index2}
				}
			}
			if v.RegKeyIndex >= 0 {
				item.Key = strings.TrimSpace(parts[v.RegKeyIndex])
			}
			for subIndex, regSubContentIndex := range v.RegSubContentIndex {
				item.SubOriginString[subIndex] = parts[regSubContentIndex]
				{
					// set index
					index1, index2 := partsIndexList[partIndex][2*regSubContentIndex], partsIndexList[partIndex][2*regSubContentIndex+1]
					if index1 >= 0 || index2 >= 0 {
						item.SubOriginIndex[subIndex] = []int{index1, index2}
					}
				}
				item.SubList[subIndex] = append(item.SubList[subIndex], c.protoBlocksFromString(item, item.SubOriginString[subIndex], v.RegSubContentTypeNames[subIndex])...)
			}
			// if item.Type.Name == GraphqlBlockExplain.Name {
			// 	log.Warn("item.SubOriginString", log.Any("item.SubOriginString", item.SubOriginString), log.Any("content", content), log.Any("parent", parent))
			// }
			// if item.Key == "gameDetails" {
			// 	log.Warn("gameDetails: ", log.Any("item.SubOriginString", item.SubOriginString), log.Any("content", content), log.Any("parent", parent))
			// }
			res = append(res, item)
		}
		if len(res) > 0 {
			return res
		}
	}
	if len(allowTypeNames) > 0 {
		lineNoComment := strings.Split(content, c.LineCommentKey)[0]
		lineNoComment = strings.TrimSpace(lineNoComment)
		if lineNoComment != "" {
			fmt.Printf("ProtoBlockFromString unknown(mustTypeName:%v allowTypeNames:%v):\n%s\n", mustTypeName, allowTypeNames, content)
		}
	}
	return nil
}

func (c *BlockParser) protoBlocksFromString(parent *Block, content string, mustTypeName []string) []*Block {
	// fmt.Println("ProtoBlocksFromString:\n" + content)
	scanner := bufio.NewScanner(strings.NewReader(content))
	res := make([]*Block, 0)
	pairCount := &PairCount{
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
		var nextLine string
		if i+1 < len(lines) {
			nextLine = lines[i+1]
		}
		_ = nextLine

		lineNoComment := strings.Split(line, c.LineCommentKey)[0]
		effectKey := pairCount.Add(lineNoComment)
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
					// log.Warn("pending line 1", log.Any("line", line), log.Any("nextLine", nextLine))
					pending = true
				}
			}
			if !pending {
				// HeadViscousPairKey check
				if effectKey != "" {
					for _, v := range c.HeadViscousPairKey {
						if v == effectKey {
							_, tail := pairKeySplit(v)
							if strings.HasSuffix(strings.TrimSpace(lineNoComment), tail) {
								pending = true
								// log.Warn("pending line 2", log.Any("line", line), log.Any("nextLine", nextLine))
								break
							}
						}
					}
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
					// log.Warn("pending line 3", log.Any("line", line), log.Any("nextLine", nextLine))
				}
			}
		}
		if !pending {
			if pairCount.IsZero() {
				res = append(res, c.protoBlockFromString(parent, currentBlockContent, mustTypeName)...)
				currentBlockContent = ""
			}
		}
	}
	if currentBlockContent != "" {
		res = append(res, c.protoBlockFromString(parent, currentBlockContent, mustTypeName)...)
	}
	return res
}
