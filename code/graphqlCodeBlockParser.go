package code

import "regexp"

var (
	GraphqlBlockTypeType         = BlockType{"type", regexp.MustCompile(`(?s)^\s*((\s*""".*?"""\n)?(\s*".*?"\n)?(extend\s+)?\s*type\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\n?(.*?[^\n]*\n?)\s*\}\n?))\s*$`), 1, 5, []int{8}, [][]string{{"type_field", "explain", "explain2"}}, []bool{true}, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeTypeField    = BlockType{"type_field", regexp.MustCompile(`(?s)^((\s*""".*?"""\n)?(\s*".*?"\n)?\s*(\w+)(\(\n?([\S\s]*?\n?)\n*?\s*\))?:\s*([\w!\[\]]+)(\s*@[^\n]*)*)\s*$`), 1, 4, []int{2, 3, 6}, [][]string{{"explain"}, {"explain2"}, {"type_field_arg"}}, []bool{true, true, false}, []string{"type"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeTypeFieldArg = BlockType{"type_field_arg", regexp.MustCompile(`(?s)(\s*""".*?"""\n)?(\s*".*?"\n)?\s*((\w+):\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)?)\s*`), 3, 4, []int{1, 2}, [][]string{{"explain"}, {"explain2"}}, []bool{true, true}, []string{"type_field"}, "", "", 0, true, nil}
	GraphqlBlockTypeInput        = BlockType{"input", regexp.MustCompile(`(?s)^(\s*""".*?"""\n)?(\s*".*?"\n)?\s*(extend\s+)?\s*input\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\n?(.*?[^\n]*\n?)\s*\}\n?)\s*$`), 0, 4, []int{7}, [][]string{{"input_field", "explain", "explain2"}}, []bool{true}, nil, "\n", "{}", 4, false, nil}
	GraphqlBlockTypeInputField   = BlockType{"input_field", regexp.MustCompile(`(?s)^((\s*""".*?"""\n)?(\s*".*?"\n)?\s*(\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)*)\s*$`), 1, 4, []int{2, 3}, [][]string{{"explain"}, {"explain2"}}, []bool{true, true}, []string{"input"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeEnum         = BlockType{"enum", regexp.MustCompile(`(?s)^\s*((extend\s+)?\s*enum\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 3, []int{6}, [][]string{{"enum_field", "explain", "explain2"}}, []bool{true}, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeEnumField    = BlockType{"enum_field", regexp.MustCompile(`(?s)^\s*((\w+)(\s*@[^\n]*)*)\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, []string{"enum"}, "\n|,", "()", 1, true, nil}
	GraphqlBlockExplain          = BlockType{"explain", regexp.MustCompile(`(?s)\s*\"\"\"\s*([^\n]*)\s*\n(([^\n]*\s*\n)*)\s*\"\"\"\s*$`), 0, 1, []int{}, [][]string{}, []bool{}, nil, "", "", 0, false, nil}
	GraphqlBlockExplain2         = BlockType{"explain2", regexp.MustCompile(`(?s)\s*\"\s*([^\n]*)\s*\"\s*$`), 0, 1, []int{}, [][]string{}, []bool{}, nil, "", "", 0, false, nil}
)

func NewGraphqlCodeBlockParser() *BlockParser {
	return &BlockParser{
		Types: []BlockType{
			GraphqlBlockTypeType,
			GraphqlBlockTypeTypeField,
			GraphqlBlockTypeTypeFieldArg,
			GraphqlBlockTypeInput,
			GraphqlBlockTypeInputField,
			GraphqlBlockTypeEnum,
			GraphqlBlockTypeEnumField,
			GraphqlBlockExplain,
			GraphqlBlockExplain2,
		},
		PairKeys:           []string{"{}", "[]", "()", `""" """`, `""`},
		HeadViscousPairKey: []string{`""" """`, `""`},
		OriginText:         []string{`""" """`},
		LineCommentKey:     "#",
		PendingLinePrefix:  "@",
	}
}
