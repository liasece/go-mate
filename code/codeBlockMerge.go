package code

import (
	"fmt"
	"strings"
)

func (b *Block) onOriginStringMerged() {
	if b.Type.RegStr == nil {
		return
	}
	// syntax
	contentReg := b.Type.RegStr
	partsList := contentReg.FindAllStringSubmatch(b.OriginString, -1)
	partsIndexList := contentReg.FindAllStringSubmatchIndex(b.OriginString, -1)
	// log.Warn("partsIndexList", log.Any("partsIndexList", partsIndexList), log.Any("partsList", partsList))
	for partIndex, parts := range partsList {
		for i, v := range parts {
			b.RegOriginStrings[i] = v
			index1, index2 := partsIndexList[partIndex][i*2], partsIndexList[partIndex][i*2+1]
			if index1 >= 0 || index2 >= 0 {
				b.RegOriginIndexes[i] = []int{index1, index2}
			}
		}
		for subIndex, regSubContentIndex := range b.Type.RegSubContentIndex {
			b.SubOriginString[subIndex] = parts[regSubContentIndex]
			{
				// set index
				index1, index2 := partsIndexList[partIndex][2*regSubContentIndex], partsIndexList[partIndex][2*regSubContentIndex+1]
				if index1 >= 0 || index2 >= 0 {
					b.SubOriginIndex[subIndex] = []int{index1, index2}
				}
			}
		}
	}
}

func (b *Block) getInParentRegOriginStringsIndex() int {
	if b.Parent == nil {
		return -1
	}
	inParentSubIndex := -1
	{
		// get this block in parent's SubList index
		for i, subList := range b.Parent.SubList {
			for _, sub := range subList {
				if sub.ID == b.ID {
					inParentSubIndex = i
					break
				}
			}
			if inParentSubIndex >= 0 {
				break
			}
		}
	}
	if inParentSubIndex < 0 {
		// not found in parent's SubList
		return -1
	}
	subRegOriginIndexesIndex := -1
	for i, v := range b.Parent.RegOriginIndexes {
		if len(v) == len(b.Parent.SubOriginIndex[inParentSubIndex]) {
			equal := true
			for j, v2 := range v {
				if b.Parent.SubOriginIndex[inParentSubIndex][j] != v2 {
					equal = false
					break
				}
			}
			if equal {
				subRegOriginIndexesIndex = i
				break
			}
		}
	}
	if subRegOriginIndexesIndex < 0 {
		// not found in parent's RegOriginIndexes
		return -1
	}
	return subRegOriginIndexesIndex
}

func (b *Block) getAddSubIntoOriginStringPos(income *Block) (int, string) {
	// get income in income's parent RegOriginStrings index
	incomeInParentRegOriginStringsIndex := income.getInParentRegOriginStringsIndex()
	if incomeInParentRegOriginStringsIndex < 0 {
		return -1, ""
	}
	{
		// find next income's parent RegOriginStrings item in b's RegOriginStrings
		resStr := income.OriginString
		incomeInParentOriginItemString := income.Parent.RegOriginStrings[incomeInParentRegOriginStringsIndex]
		for i := incomeInParentRegOriginStringsIndex + 1; i < len(income.Parent.RegOriginStrings); i++ {
			if income.Parent.RegOriginStrings[i] == "" {
				continue
			}
			for j, v := range b.RegOriginStrings {
				if v == income.Parent.RegOriginStrings[i] {
					// find ok
					pos := b.RegOriginIndexes[j][0]
					for pos > 0 && (b.OriginString[pos-1] == ' ' || b.OriginString[pos-1] == '\t' || b.OriginString[pos-1] == '\n') {
						pos--
					}
					return pos, resStr
				}
			}
			if strings.Contains(income.Parent.RegOriginStrings[i], incomeInParentOriginItemString) {
				resStr = strings.Replace(income.Parent.RegOriginStrings[i], incomeInParentOriginItemString, resStr, 1)
				incomeInParentOriginItemString = income.Parent.RegOriginStrings[i]
			}
		}
	}
	{
		// find pre income's parent RegOriginStrings item in b's RegOriginStrings
		resStr := income.OriginString
		incomeInParentOriginItemString := income.Parent.RegOriginStrings[incomeInParentRegOriginStringsIndex]
		for i := incomeInParentRegOriginStringsIndex - 1; i >= 0; i-- {
			if income.Parent.RegOriginStrings[i] == "" {
				continue
			}
			for j, v := range b.RegOriginStrings {
				if v == income.Parent.RegOriginStrings[i] {
					// find ok
					return b.RegOriginIndexes[j][1], resStr
				}
			}
			if strings.Contains(income.Parent.RegOriginStrings[i], incomeInParentOriginItemString) {
				resStr = strings.Replace(income.Parent.RegOriginStrings[i], incomeInParentOriginItemString, resStr, 1)
				incomeInParentOriginItemString = income.Parent.RegOriginStrings[i]
			}
		}
	}
	return -1, ""
}

func (b *Block) addSub(targetSubLevel int, income *Block) {
	if targetSubLevel <= 0 {
		targetSubLevel = 1
	}
	if len(b.SubList) < targetSubLevel {
		panic(fmt.Sprintf("addSub error: %s(%s) targetSubLevel: %d", b.Key, b.Type.Name, targetSubLevel))
		return
	}
	targetSubIndex := targetSubLevel - 1

	joinString := b.getSubJoinString(targetSubIndex)
	tailString := ""
	tabStr := ""
	if joinString == "\n" {
		tailString = "\n"
		tabStr = b.getSubTabStr(joinString)
	}
	// fmt.Println("addSub: " + b.Key + " joinString: ```" + joinString + "```\ntailString: ```" + tailString + "```")
	newSubOriginString := income.OriginString
	if b.SubOriginString[targetSubIndex] != "" {
		newSubOriginString = fmt.Sprintf("%s%s%s%s", b.SubOriginString[targetSubIndex], joinString, tabStr, income.OriginString)
	}
	income.Parent = b
	b.SubList[targetSubIndex] = append(b.SubList[targetSubIndex], income)
	if b.SubOriginString[targetSubIndex] == "" {
		if insertPos, insertStr := b.getAddSubIntoOriginStringPos(income); insertPos >= 0 {
			// new subs
			myOldOriginString := b.OriginString
			b.OriginString = fmt.Sprintf("%s%s%s%s%s", b.OriginString[:insertPos], insertStr, tailString, tabStr, b.OriginString[insertPos:])
			b.onOriginStringMerged()
			if b.Parent != nil {
				newSubOriginString = strings.Replace(b.Parent.SubOriginString[targetSubIndex], myOldOriginString, b.OriginString, 1)
			}
		} else if b.Type.SubWarpChar != "" && b.Type.RegSubWarpContentIndex > 0 {
			// get sub warp content
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
			b.SubOriginString[targetSubIndex] = strings.Trim(newSubOriginString, " \t\r\n")
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
				b.OriginString = fmt.Sprintf("%s%s%s%s%s%s%s", b.OriginString[:insertPos], subWarpCharHead, tabStr, b.SubOriginString[targetSubIndex], tailString, subWarpCharTail, b.OriginString[insertPos:])
			} else {
				b.OriginString = fmt.Sprintf("%s%s%s%s%s", b.OriginString[:insertPos], tabStr, b.SubOriginString[targetSubIndex], tailString, b.OriginString[insertPos:])
			}
			b.onOriginStringMerged()

			if b.Parent != nil {
				newSubOriginString = strings.Replace(b.Parent.SubOriginString[targetSubIndex], myOldOriginString, b.OriginString, 1)
			}
		}
		b = b.Parent
	}

	{
		// update parent content
		for b != nil {
			myOldSubOriginString := b.SubOriginString[targetSubIndex]
			myOldOriginString := b.OriginString

			// update sub origin
			b.SubOriginString[targetSubIndex] = newSubOriginString

			// update origin
			b.OriginString = strings.Replace(b.OriginString, myOldSubOriginString, b.SubOriginString[targetSubIndex], 1)
			// fmt.Println(b.Key + " new OriginString: ```" + b.OriginString + "```\nmyOldOriginString: ```" + myOldOriginString + "```" + "```\nreplaceFrom: ```" + myOldSubOriginString + "```" + "\nreplaceTo: ```" + b.SubOriginString + "```\n")
			b.onOriginStringMerged()

			if b.Parent != nil {
				newSubOriginString = strings.Replace(b.Parent.SubOriginString[targetSubIndex], myOldOriginString, b.OriginString, 1)
			}
			b = b.Parent
		}
	}
}

// subLevel = subIndex + 1
func (b *Block) Merge(targetSubLevel int, income *Block) *Block {
	subLevel, findBlock := b.Find(income.Type, income.Key, targetSubLevel)
	if findBlock == nil {
		// fmt.Println("Merge not exists:" + income.Key + "(b: " + b.Key + "(" + b.Type.Name + ")" + ")" + "(income: " + income.Key + "(" + income.Type.Name + ")" + ")")
		// append to current block
		b.addSub(subLevel, income)
	} else {
		for subIndex, subs := range income.SubList {
			for _, v := range subs {
				mergeAble := false
				if targetSubLevel == 0 {
					mergeAble = true
				} else if (subLevel <= 0 || findBlock.Type.SubMergeType[subLevel-1]) && income.Type.SubMergeType[subIndex] {
					mergeAble = true
				}
				if mergeAble {
					findBlock.Merge(subIndex+1, v)
				}
			}
		}
	}

	return b
}
