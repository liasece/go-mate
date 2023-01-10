package code

import (
	"regexp"
	"strings"
)

type Block struct {
	ID           string
	BlockParser  *BlockParser `json:"-"`
	Parent       *Block       `json:"-"`
	Key          string
	Type         BlockType
	OriginString string

	RegOriginStrings []string
	RegOriginIndexes [][]int

	SubOriginString []string
	SubOriginIndex  [][]int

	SubList [][]*Block
}

func (b *Block) Clone() *Block {
	return b.clone(b.Parent)
}

func (b *Block) clone(parent *Block) *Block {
	var regOriginStrings []string
	if b.RegOriginStrings != nil {
		regOriginStrings = make([]string, 0)
		regOriginStrings = append(regOriginStrings, b.RegOriginStrings...)
	}
	var regOriginIndexes [][]int
	if b.RegOriginIndexes != nil {
		regOriginIndexes = make([][]int, len(b.RegOriginIndexes))
		for i, v := range b.RegOriginIndexes {
			if v != nil {
				regOriginIndexes[i] = make([]int, 0)
				regOriginIndexes[i] = append(regOriginIndexes[i], v...)
			}
		}
	}
	var subOriginString []string
	if b.SubOriginString != nil {
		subOriginString = make([]string, 0)
		subOriginString = append(subOriginString, b.SubOriginString...)
	}
	var subOriginIndex [][]int
	if b.SubOriginIndex != nil {
		subOriginIndex = make([][]int, len(b.SubOriginIndex))
		for i, v := range b.SubOriginIndex {
			if v != nil {
				subOriginIndex[i] = make([]int, 0)
				subOriginIndex[i] = append(subOriginIndex[i], v...)
			}
		}
	}
	res := &Block{
		ID:               b.ID,
		BlockParser:      b.BlockParser,
		Parent:           parent,
		Key:              b.Key,
		Type:             b.Type,
		OriginString:     b.OriginString,
		RegOriginStrings: regOriginStrings,
		RegOriginIndexes: regOriginIndexes,
		SubOriginString:  subOriginString,
		SubOriginIndex:   subOriginIndex,
		SubList:          nil,
	}
	var subList [][]*Block
	if b.SubList != nil {
		subList = make([][]*Block, len(b.SubList))
		for i, v := range b.SubList {
			if v != nil {
				subList[i] = make([]*Block, len(v))
				for j, v2 := range v {
					subList[i][j] = v2.clone(res)
				}
			}
		}
	}
	res.SubList = subList
	return res
}

// subLevel = subIndex + 1
func (b *Block) Find(typ BlockType, key string, targetSubLevel int) (findBlock *Block) {
	if targetSubLevel == 0 {
		bKey := b.Key
		findKey := key
		if b.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
			findKey = strings.ToLower(findKey)
			bKey = strings.ToLower(bKey)
		}
		if b.Type.Name == typ.Name && bKey == findKey {
			return b
		}
	} else {
		for subIndex, subs := range b.SubList {
			if subIndex+1 != targetSubLevel {
				continue
			}
			for _, v := range subs {
				vKey := v.Key
				findKey := key
				if v.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
					vKey = strings.ToLower(vKey)
					findKey = strings.ToLower(findKey)
				}
				if v.Type.Name == typ.Name && vKey == findKey {
					return v
				}
			}
		}
	}
	return nil
}

func (b *Block) getFirstSubBlock() *Block {
	for _, v := range b.SubList {
		if len(v) > 0 {
			return v[0]
		}
	}
	return nil
}

func (b *Block) allSubBlocks() []*Block {
	blocks := []*Block{}
	for _, v := range b.SubList {
		blocks = append(blocks, v...)
	}
	return blocks
}

func (b *Block) getSubJoinString(targetSubIndex int) string {
	if len(b.SubList[targetSubIndex]) > 1 {
		// try get from origin sub string
		leftString := strings.TrimSpace(b.SubList[targetSubIndex][0].OriginString)
		rightString := strings.TrimSpace(b.SubList[targetSubIndex][1].OriginString)
		firstSubAndJoinString := strings.Split(b.SubOriginString[targetSubIndex], rightString)[0]
		joinString := strings.Split(firstSubAndJoinString, leftString)
		if len(joinString) > 1 {
			return joinString[len(joinString)-1]
		}
	}
	var firstSubAndJoinString string
	if len(b.SubList[targetSubIndex]) == 1 {
		// try get from origin sub string
		rightString := strings.TrimSpace(b.SubList[targetSubIndex][0].OriginString)
		firstSubAndJoinString = strings.Split(b.SubList[targetSubIndex][0].OriginString, rightString)[0]
	}
	if b.Type.SubsSeparator != "" {
		oldSub := ""
		subs := strings.Split(b.Type.SubsSeparator, "|")
		for _, sub := range subs {
			var subSliceTail string
			if len(b.SubList[targetSubIndex]) > 0 {
				subSlice := strings.Split(b.SubOriginString[targetSubIndex], b.SubList[targetSubIndex][0].OriginString)
				if len(subSlice) > 1 {
					subSliceTail = subSlice[1]
				}
			}
			if len(b.SubList[targetSubIndex]) > 1 && strings.HasPrefix(subSliceTail, sub) {
				oldSub = sub
				break
			}
		}
		if oldSub == "" {
			oldSub = subs[0]
		}
		if strings.Trim(oldSub, " \n\r\t") != "" {
			oldSub += " "
		}
		return oldSub + firstSubAndJoinString
	}
	return ""
}

func (b *Block) getSubTabStr(originTail string) string {
	str := ""
	if firstSubBlock := b.getFirstSubBlock(); firstSubBlock != nil {
		contentReg := regexp.MustCompile(`(?s)(` + originTail + `)?(\s*)` + regexp.QuoteMeta(firstSubBlock.OriginString))
		partsList := contentReg.FindAllStringSubmatch(b.OriginString, -1)
		if len(partsList) > 0 {
			parts := partsList[0]
			if len(parts) > 1 {
				str += parts[2]
			}
		}
	} else if b.Parent != nil {
		for _, v := range b.Parent.allSubBlocks() {
			if v == b {
				continue
			}
			if firstSubBlock := v.getFirstSubBlock(); firstSubBlock != nil {
				contentReg := regexp.MustCompile(`(?s)(` + originTail + `)?(\s*)` + regexp.QuoteMeta(firstSubBlock.OriginString))
				partsList := contentReg.FindAllStringSubmatch(v.OriginString, -1)
				if len(partsList) > 0 {
					parts := partsList[0]
					if len(parts) > 1 {
						str += parts[2]
						break
					}
				}
			}
		}
	}
	str = strings.Replace(str, originTail, "", 1)
	if b.Parent != nil {
		str = b.Parent.getSubTabStr(originTail) + str
	}
	return str
}
