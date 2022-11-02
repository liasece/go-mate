package code

import (
	"fmt"
	"regexp"
	"strings"
)

type CodeBlock struct {
	Parent          *CodeBlock `json:"-"`
	Key             string
	Type            CodeBlockType
	OriginString    string
	SubOriginString string
	SubList         []*CodeBlock
}

func (b *CodeBlock) Find(typ CodeBlockType, key string) *CodeBlock {
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

func (b *CodeBlock) getSubJoinString(income *CodeBlock) string {
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

func (b *CodeBlock) addSub(income *CodeBlock) {
	joinString := b.getSubJoinString(income)
	tailString := ""
	if joinString == "\n" {
		tailString = "\n"
	}
	// fmt.Println("addSub: " + b.Key + " joinString: ```" + joinString + "```\ntailString: ```" + tailString + "```")
	newSubOriginString := income.OriginString
	if b.SubOriginString != "" {
		newSubOriginString = fmt.Sprintf("%s%s%s", b.SubOriginString, joinString, income.OriginString)
	}
	income.Parent = b
	b.SubList = append(b.SubList, income)
	if b.SubOriginString == "" && b.Type.SubWarpChar != "" && b.Type.RegSubWarpContentIndex > 0 {
		// new subs
		myOldOriginString := b.OriginString
		// insert first sub origin string
		b.SubOriginString = strings.Trim(newSubOriginString, " \t\r\n")
		insertPos := -1
		newBlock := false
		{
			// find insert pos
			contentReg := regexp.MustCompile(b.Type.RegStr)
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
			b.OriginString = fmt.Sprintf("%s%s%s%s%s%s", b.OriginString[:insertPos], b.Type.SubWarpChar[:1], b.SubOriginString, tailString, b.Type.SubWarpChar[1:], b.OriginString[insertPos:])
		} else {
			b.OriginString = fmt.Sprintf("%s%s%s%s", b.OriginString[:insertPos], b.SubOriginString, tailString, b.OriginString[insertPos:])
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

func (b *CodeBlock) Merge(income *CodeBlock) *CodeBlock {
	exists := b.Find(income.Type, income.Key)
	if exists == nil {
		// fmt.Println("Merge not exists:" + income.Key + "(b: " + b.Key + "(" + b.Type.Name + ")" + ")" + "(income: " + income.Key + "(" + income.Type.Name + ")" + ")")
		// append to current block
		b.addSub(income)
	} else {
		// fmt.Println("Merge exists:" + income.Key + "(" + exists.Type.Name + ")" + " (MergeAble: " + fmt.Sprint(exists.Type.MergeAble) + ")")
		if exists.Type.MergeAble && income.Type.MergeAble {
			for _, v := range income.SubList {
				exists.Merge(v)
			}
		}
		// switch exists.Type.Name {
		// case ProtoBlockTypeService.Name, ProtoBlockTypeMessage.Name, ProtoBlockTypeOption.Name, ProtoBlockTypeNone.Name, ProtoBlockTypeMessageField.Name:
		// 	for _, v := range income.SubList {
		// 		exists.Merge(v)
		// 	}
		// }
	}
	return b
}
