package code

import (
	"fmt"
	"regexp"
	"strings"
)

type Block struct {
	Parent          *Block `json:"-"`
	Key             string
	Type            BlockType
	OriginString    string
	SubOriginString string
	SubList         []*Block
}

func (b *Block) Find(typ BlockType, key string) *Block {
	bKey := b.Key
	if b.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
		key = strings.ToLower(key)
		bKey = strings.ToLower(bKey)
	}
	if b.Type.Name == typ.Name && bKey == key {
		return b
	}
	for _, v := range b.SubList {
		vKey := v.Key
		if v.Type.KeyCaseIgnored || typ.KeyCaseIgnored {
			vKey = strings.ToLower(vKey)
		}
		if v.Type.Name == typ.Name && vKey == key {
			return v
		}
	}
	return nil
}

func (b *Block) getSubJoinString() string {
	if b.Type.SubsSeparator != "" {
		oldSub := ""
		subs := strings.Split(b.Type.SubsSeparator, "|")
		for _, sub := range subs {
			if len(b.SubList) > 1 && strings.HasPrefix(strings.Split(b.SubOriginString, b.SubList[0].OriginString)[1], sub) {
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
	if len(b.SubList) > 0 {
		contentReg := regexp.MustCompile(`(?s)(` + originTail + `)?(\s*)` + regexp.QuoteMeta(b.SubList[0].OriginString))
		partsList := contentReg.FindAllStringSubmatch(b.OriginString, -1)
		if len(partsList) > 0 {
			parts := partsList[0]
			if len(parts) > 1 {
				str += parts[2]
			}
		}
	} else if b.Parent != nil {
		for _, v := range b.Parent.SubList {
			if v == b {
				continue
			}
			if len(v.SubList) > 0 {
				contentReg := regexp.MustCompile(`(?s)(` + originTail + `)?(\s*)` + regexp.QuoteMeta(v.SubList[0].OriginString))
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

func (b *Block) addSub(income *Block) {
	joinString := b.getSubJoinString()
	tailString := ""
	tabStr := ""
	if joinString == "\n" {
		tailString = "\n"
		tabStr = b.getSubTabStr(joinString)
	}
	// fmt.Println("addSub: " + b.Key + " joinString: ```" + joinString + "```\ntailString: ```" + tailString + "```")
	newSubOriginString := income.OriginString
	if b.SubOriginString != "" {
		newSubOriginString = fmt.Sprintf("%s%s%s%s", b.SubOriginString, joinString, tabStr, income.OriginString)
	}
	income.Parent = b
	b.SubList = append(b.SubList, income)
	if b.SubOriginString == "" && b.Type.SubWarpChar != "" && b.Type.RegSubWarpContentIndex > 0 {
		subWarpCharHead := ""
		subWarpCharTail := ""
		if len(b.Type.SubWarpChar) == 2 {
			subWarpCharHead = b.Type.SubWarpChar[:1]
			subWarpCharTail = b.Type.SubWarpChar[1:]
		} else if ss := strings.Split(b.Type.SubWarpChar, "|"); len(ss) == 2 {
			subWarpCharHead = ss[0]
			subWarpCharTail = ss[1]
		} else {
			subWarpCharHead = b.Type.SubWarpChar
			subWarpCharTail = b.Type.SubWarpChar
		}
		// new subs
		myOldOriginString := b.OriginString
		// insert first sub origin string
		b.SubOriginString = strings.Trim(newSubOriginString, " \t\r\n")
		insertPos := -1
		newBlock := false
		{
			// find insert pos
			contentReg := b.Type.RegStr
			indexes := contentReg.FindStringSubmatchIndex(b.OriginString)
			matchIndex := b.Type.RegSubWarpContentIndex
			if indexes[matchIndex*2+1] >= 0 {
				insertPos = indexes[matchIndex*2+1] - 1
			} else {
				for matchIndex > 0 {
					matchIndex--
					if indexes[matchIndex*2] >= 0 {
						insertPos = indexes[matchIndex*2+1]
						newBlock = true
						break
					}
				}
			}
		}
		if insertPos < 0 {
			panic("addSub to empty block: " + income.Key)
		}
		if newBlock {
			b.OriginString = fmt.Sprintf("%s%s%s%s%s%s%s", b.OriginString[:insertPos], subWarpCharHead, tabStr, b.SubOriginString, tailString, subWarpCharTail, b.OriginString[insertPos:])
		} else {
			b.OriginString = fmt.Sprintf("%s%s%s%s%s", b.OriginString[:insertPos], tabStr, b.SubOriginString, tailString, b.OriginString[insertPos:])
		}

		if b.Parent != nil {
			newSubOriginString = strings.Replace(b.Parent.SubOriginString, myOldOriginString, b.OriginString, 1)
		}
		b = b.Parent
	}
	for b != nil {
		myOldSubOriginString := b.SubOriginString
		myOldOriginString := b.OriginString

		// update sub origin
		b.SubOriginString = newSubOriginString

		// update origin
		b.OriginString = strings.Replace(b.OriginString, myOldSubOriginString, b.SubOriginString, 1)
		// fmt.Println(b.Key + " new OriginString: ```" + b.OriginString + "```\nmyOldOriginString: ```" + myOldOriginString + "```" + "```\nreplaceFrom: ```" + myOldSubOriginString + "```" + "\nreplaceTo: ```" + b.SubOriginString + "```\n")

		if b.Parent != nil {
			newSubOriginString = strings.Replace(b.Parent.SubOriginString, myOldOriginString, b.OriginString, 1)
		}
		b = b.Parent
	}
}

func (b *Block) Merge(income *Block) *Block {
	exists := b.Find(income.Type, income.Key)
	if exists == nil {
		// fmt.Println("Merge not exists:" + income.Key + "(b: " + b.Key + "(" + b.Type.Name + ")" + ")" + "(income: " + income.Key + "(" + income.Type.Name + ")" + ")")
		// append to current block
		b.addSub(income)
	} else if exists.Type.MergeAble && income.Type.MergeAble {
		for _, v := range income.SubList {
			exists.Merge(v)
		}
	}
	{
		// add end line
		if b.OriginString != "" && b.OriginString[len(b.OriginString)-1] != '\n' {
			b.OriginString += "\n"
		}
	}
	return b
}
