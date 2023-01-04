package code

import "regexp"

var (
	GraphqlBlockTypeType         = BlockType{"type", true, regexp.MustCompile(`(?s)^\s*((extend\s+)?\s*type\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 3, 6, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeTypeField    = BlockType{"type_field", false, regexp.MustCompile(`(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*@[^\n]*)*)\s*$`), 1, 2, 4, []string{"type"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeTypeFieldArg = BlockType{"type_field_arg", false, regexp.MustCompile(`(?s)\s*((\w+):\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)?)\s*`), 1, 2, -1, []string{"type_field"}, "", "", 0, true, nil}
	GraphqlBlockExplain          = BlockType{"explain", true, regexp.MustCompile(`(?s)\s*\"\"\"\s*([^\n]*\n)\s*(([^\n]*\s*)*)\s*\"\"\"\s*$`), 0, 1, -1, nil, "", "", 0, false, nil}
	GraphqlBlockExplain2         = BlockType{"explain2", true, regexp.MustCompile(`(?s)\s*\"([^\n]*)\"\s*$`), 0, 1, -1, nil, "", "", 0, false, nil}
	GraphqlBlockTypeInput        = BlockType{"input", true, regexp.MustCompile(`(?s)^\s*(extend\s+)?\s*input\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\})\s*$`), 0, 2, 5, nil, "\n", "{}", 4, false, nil}
	GraphqlBlockTypeInputField   = BlockType{"input_field", true, regexp.MustCompile(`(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)*)\s*$`), 1, 2, 4, []string{"input"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeEnum         = BlockType{"enum", true, regexp.MustCompile(`(?s)^\s*((extend\s+)?\s*enum\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 3, 6, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeEnumField    = BlockType{"enum_field", true, regexp.MustCompile(`(?s)^\s*((\w+)(\s*@[^\n]*)*)\s*$`), 1, 2, -1, []string{"enum"}, "\n|,", "()", 1, true, nil}
)

func NewGraphqlCodeBlockParser() *BlockParser {
	return &BlockParser{
		Types: []BlockType{
			GraphqlBlockTypeType,
			GraphqlBlockTypeTypeField,
			GraphqlBlockTypeTypeFieldArg,
			GraphqlBlockExplain,
			GraphqlBlockExplain2,
			GraphqlBlockTypeInput,
			GraphqlBlockTypeInputField,
			GraphqlBlockTypeEnum,
			GraphqlBlockTypeEnumField,
		},
		PairKeys: []string{"{}", "[]", "()", `""" """`, `""`},
		// HeadViscousPairKey: []string{`""" """`, `""`},
		OriginText:        []string{`""" """`},
		LineCommentKey:    "#",
		PendingLinePrefix: "@",
	}
}
