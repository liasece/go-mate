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
	if b.Type.Name == typ.Name && b.Key == key {
		return b
	}
	for _, v := range b.SubList {
		if v.Type.Name == typ.Name && v.Key == key {
			return v
		}
	}
	return nil
}

func (b *CodeBlock) addSub(income *CodeBlock) {
	newSubOriginString := fmt.Sprintf("%s%s", b.SubOriginString, income.OriginString)
	if income.Type.SubsSeparator != "" && len(b.SubList) > 0 {
		if !strings.HasSuffix(strings.TrimSpace(b.SubOriginString), income.Type.SubsSeparator) {
			newSubOriginString = fmt.Sprintf("%s%s %s", b.SubOriginString, income.Type.SubsSeparator, income.OriginString)
		}
	}
	income.Parent = b
	b.SubList = append(b.SubList, income)
	if b.SubOriginString == "" {
		myOldOriginString := b.OriginString
		// insert first sub origin string
		if b.Type.SubWarpChar == "" || b.Type.RegSubWarpContentIndex <= 0 {
			panic("addSub to empty block: " + income.Key + "(" + b.Type.Name + ")")
		}
		b.SubOriginString = newSubOriginString
		insertPos := -1
		{
			// find insert pos
			contentReg := regexp.MustCompile(b.Type.RegStr)
			indexes := contentReg.FindStringSubmatchIndex(b.OriginString)
			matchIndex := b.Type.RegSubWarpContentIndex
			if indexes[matchIndex*2] >= 0 {
				insertPos = indexes[matchIndex*2]
			} else {
				for matchIndex > 0 {
					matchIndex--
					if indexes[matchIndex*2] >= 0 {
						insertPos = indexes[matchIndex*2+1]
						break
					}
				}
			}
		}
		if insertPos < 0 {
			panic("addSub to empty block: " + income.Key)
		}
		b.OriginString = fmt.Sprintf("%s%s%s%s%s", b.OriginString[:insertPos], b.Type.SubWarpChar[:1], b.SubOriginString, b.Type.SubWarpChar[1:], b.OriginString[insertPos:])

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
		// fmt.Println(b.Key + " new OriginString: " + b.OriginString + " (replaceFrom: " + myOldSubOriginString + ")" + " (replaceTo: " + b.SubOriginString + ")")

		if b.Parent != nil {
			newSubOriginString = strings.Replace(b.Parent.SubOriginString, myOldOriginString, b.OriginString, 1)
		}
		b = b.Parent
	}
}

func (b *CodeBlock) Merge(income *CodeBlock) *CodeBlock {
	exists := b.Find(income.Type, income.Key)
	if exists == nil {
		// fmt.Println("Merge not exists:" + income.Key)
		// append to current block
		b.addSub(income)
	} else {
		// fmt.Println("Merge exists:" + income.Key)
		// if exists.Type.MergeAble {
		// 	for _, v := range income.SubList {
		// 		exists.Merge(v)
		// 	}
		// }
		switch exists.Type.Name {
		case ProtoBlockTypeService.Name, ProtoBlockTypeMessage.Name, ProtoBlockTypeOption.Name, ProtoBlockTypeNone.Name, ProtoBlockTypeMessageField.Name:
			for _, v := range income.SubList {
				exists.Merge(v)
			}
		}
	}
	return b
}
