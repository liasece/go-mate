package code

import (
	"fmt"
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
	if len(b.SubList) == 0 {
		panic("addSub to empty block: " + income.Key)
	}
	income.Parent = b
	b.SubList = append(b.SubList, income)
	newSubOriginString := fmt.Sprintf("%s%s", b.SubOriginString, income.OriginString)
	switch income.Type.Name {
	case ProtoBlockTypeOptionItem.Name:
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
