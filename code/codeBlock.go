package code

import (
	"regexp"
	"strings"
)

type Block struct {
	ID               string
	Parent           *Block `json:"-"`
	Key              string
	Type             BlockType
	OriginString     string
	RegOriginStrings []string
	RegOriginIndexes [][]int
	SubOriginString  []string
	SubOriginIndex   [][]int
	SubList          [][]*Block
}

// subLevel = subIndex + 1
func (b *Block) Find(typ BlockType, key string, targetSubLevel int) (subLevel int, findBlock *Block) {
	bKey := b.Key
	if b.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
		key = strings.ToLower(key)
		bKey = strings.ToLower(bKey)
	}
	if b.Type.Name == typ.Name && bKey == key && targetSubLevel <= 0 {
		return 0, b
	}
	for subIndex, subs := range b.SubList {
		if targetSubLevel > 0 && subIndex+1 != targetSubLevel {
			continue
		}
		for _, v := range subs {
			vKey := v.Key
			if v.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
				vKey = strings.ToLower(vKey)
			}
			if v.Type.Name == typ.Name && vKey == key {
				return subIndex + 1, v
			}
		}
	}
	return 0, nil
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
		firstSubAndJoinString := strings.Split(b.SubOriginString[targetSubIndex], b.SubList[targetSubIndex][1].OriginString)[0]
		joinString := strings.Split(firstSubAndJoinString, b.SubList[targetSubIndex][0].OriginString)
		if len(joinString) > 1 {
			return joinString[len(joinString)-1]
		}
	}
	if b.Type.SubsSeparator != "" {
		oldSub := ""
		subs := strings.Split(b.Type.SubsSeparator, "|")
		for _, sub := range subs {
			if len(b.SubList[targetSubIndex]) > 1 && strings.HasPrefix(strings.Split(b.SubOriginString[targetSubIndex], b.SubList[targetSubIndex][0].OriginString)[1], sub) {
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
		return oldSub
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
