package code

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// MergeConfig is the config for merge code block.
type MergeConfig struct {
	Append           bool
	ReplaceBlockType []string
}

type BlockType struct {
	Name                   string
	RegStr                 *regexp.Regexp
	RegOriginIndex         int
	RegKeyIndex            int
	RegSubContentIndex     []int      // sub content in reg index list
	RegSubContentTypeNames [][]string // sub content type name list
	SubMergeType           []*MergeConfig
	ParentNames            []string
	SubsSeparator          string // the sub body separator character, like "," or ";"
	SubWarpChar            string // the sub body warp character, like "()" "{}" "[]" "((|))" `"""|"""`
	RegSubWarpContentIndex int
	KeyCaseIgnored         bool     // ABc == abc
	SubTailChar            []string // the sub body tail character, like "}" ";" "))" ","
}

var _ json.Marshaler = (*BlockType)(nil)

func (t *BlockType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", t.Name)), nil
}
