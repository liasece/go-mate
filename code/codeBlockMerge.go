package code

import (
	"fmt"
	"strings"

	"github.com/liasece/log"
)

func (b *Block) onOriginStringMerged() {
	if b.Type.RegStr == nil {
		return
	}
	// syntax
	partsList := b.Type.RegStr.FindAllStringSubmatch(b.OriginString, -1)
	partsIndexList := b.Type.RegStr.FindAllStringSubmatchIndex(b.OriginString, -1)
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
			subs := b.BlockParser.protoBlocksFromString(b, b.SubOriginString[subIndex], nil)
			for _, sub := range b.SubList[subIndex] {
				var newSub *Block
				for _, v := range subs {
					if v.Type.Name == sub.Type.Name && v.Key == sub.Key {
						newSub = v
						break
					}
				}
				if newSub != nil {
					sub.OriginString = newSub.OriginString
					sub.onOriginStringMerged()
				} else {
					log.Warn("onOriginStringMerged not found int new sub list", log.String("subType", sub.Type.Name), log.String("subKey", sub.Key))
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

func (b *Block) getAddSubIntoOriginStringPos(income *Block, targetSubIndex int) (int, string) {
	// get income in income's parent RegOriginStrings index
	incomeInParentRegOriginStringsIndex := income.getInParentRegOriginStringsIndex()
	if incomeInParentRegOriginStringsIndex < 0 {
		return -1, ""
	}

	{
		// like merge {a = 1;} into {}
		if strings.Replace(income.Parent.OriginString, income.OriginString, "", 1) == b.OriginString {
			return strings.Index(income.Parent.OriginString, income.OriginString), income.OriginString
		}
	}

	insertStringTailString := ""
	{
		// get sub list join string
		subListJoinString := b.getSubJoinString(targetSubIndex)
		if subListJoinString == "" && income.Parent != nil {
			incomeParentSubJoinString := income.Parent.getSubJoinString(income.getInParentSubLevel() - 1)
			if incomeParentSubJoinString != "" {
				subListJoinString = incomeParentSubJoinString
			}
		}
		if strings.Contains(subListJoinString, "\n") {
			insertStringTailString = "\n"
		}
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
					return pos, resStr + insertStringTailString
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
					return b.RegOriginIndexes[j][1], resStr + insertStringTailString
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

func (b *Block) getInParentSubLevel() int {
	if b.Parent == nil {
		return -1
	}
	inParentSubLevel := -1
	if b.ID == b.Parent.ID {
		inParentSubLevel = 0
	} else {
		// get this block in parent's SubList index
		for i, subList := range b.Parent.SubList {
			for _, sub := range subList {
				if sub.ID == b.ID {
					inParentSubLevel = i + 1
					break
				}
			}
			if inParentSubLevel > 0 {
				break
			}
		}
	}
	if inParentSubLevel < 0 {
		// not found in parent's SubList
		return -1
	}
	return inParentSubLevel
}

// update this block's RegOriginStrings, if this block has parent, update parent's info also
func (b *Block) updateSubOriginString(targetSubIndex int, newSubOriginString string) {
	myOldSubOriginString := b.SubOriginString[targetSubIndex]
	myOldOriginString := b.OriginString

	// update this block's OriginString
	// update sub origin
	b.SubOriginString[targetSubIndex] = newSubOriginString

	// update origin
	b.OriginString = strings.Replace(b.OriginString, myOldSubOriginString, newSubOriginString, 1)
	b.onOriginStringMerged()

	// update this block's parent's info
	if b.Parent != nil {
		inParentSubIndex := b.getInParentSubLevel() - 1
		parentNewOriginString := strings.Replace(b.Parent.SubOriginString[inParentSubIndex], myOldOriginString, b.OriginString, 1)
		b.Parent.updateSubOriginString(inParentSubIndex, parentNewOriginString)
	}
}

func (b *Block) addSub(targetSubLevel int, income *Block) {
	if targetSubLevel <= 0 {
		targetSubLevel = 1
	}
	if len(b.SubList) < targetSubLevel {
		panic(fmt.Sprintf("addSub error: %s(%s) targetSubLevel: %d", b.Key, b.Type.Name, targetSubLevel))
	}
	targetSubIndex := targetSubLevel - 1

	// get sub list join string
	subListJoinString := b.getSubJoinString(targetSubIndex)
	if subListJoinString == "" && income.Parent != nil {
		incomeParentSubJoinString := income.Parent.getSubJoinString(income.getInParentSubLevel() - 1)
		if incomeParentSubJoinString != "" {
			subListJoinString = incomeParentSubJoinString
		}
	}

	{
		newSubBlock := income.Clone()
		newSubBlock.Parent = b
		b.SubList[targetSubIndex] = append(b.SubList[targetSubIndex], newSubBlock)
	}
	if b.SubOriginString[targetSubIndex] != "" {
		// append to sub list
		newSubOriginString := fmt.Sprintf("%s%s%s", strings.TrimRightFunc(b.SubOriginString[targetSubIndex], func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }), subListJoinString, strings.TrimLeftFunc(income.OriginString, func(r rune) bool { return r == ' ' || r == '\t' || r == '\n' }))
		// log.Error("append to sub list", log.Any("bKey", b.Key), log.Any("bType", b.Type.Name), log.Any("targetSubLevel", targetSubLevel), log.Any("incomeKey", income.Key), log.Any("incomeType", income.Type.Name), log.Any("subListJoinString", subListJoinString), log.Any("newSubOriginString", newSubOriginString))
		b.updateSubOriginString(targetSubIndex, newSubOriginString)
	} else if insertPos, insertStr := b.getAddSubIntoOriginStringPos(income, targetSubIndex); insertPos >= 0 {
		// new sub list
		myOldOriginString := b.OriginString
		b.OriginString = fmt.Sprintf("%s%s%s", b.OriginString[:insertPos], insertStr, b.OriginString[insertPos:])
		// log.Error("new sub list", log.Any("bKey", b.Key), log.Any("bType", b.Type.Name), log.Any("targetSubLevel", targetSubLevel), log.Any("incomeKey", income.Key), log.Any("incomeType", income.Type.Name), log.Any("newBOriginString", b.OriginString), log.Any("insertStr", insertStr))
		b.onOriginStringMerged()
		if b.Parent != nil {
			newSubOriginString := strings.Replace(b.Parent.SubOriginString[targetSubIndex], myOldOriginString, b.OriginString, 1)
			b.Parent.updateSubOriginString(b.getInParentSubLevel()-1, newSubOriginString)
		}
	} else {
		log.Panic("can't add sub list to this block", log.Any("bKey", b.Key), log.Any("bType", b.Type.Name), log.Any("targetSubLevel", targetSubLevel), log.Any("incomeKey", income.Key), log.Any("incomeType", income.Type.Name))
	}
}

// subLevel = subIndex + 1
func (b *Block) Merge(targetSubLevel int, income *Block) *Block {
	// log.Warn("in Merge", log.Any("bKey", b.Key), log.Any("targetSubLevel", targetSubLevel), log.Any("incomeKey", income.Key))
	subLevel, findBlock := b.Find(income.Type, income.Key, targetSubLevel)
	if findBlock == nil {
		if b.Type.Name != income.Type.Name {
			subLevel = targetSubLevel
		}

		// log.Warn("Merge not exists, add sub", log.Any("bKey", b.Key), log.Any("bType", b.Type.Name), log.Any("targetSubLevel", targetSubLevel), log.Any("incomeKey", income.Key), log.Any("incomeType", income.Type.Name), log.Any("subLevel", subLevel))
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
