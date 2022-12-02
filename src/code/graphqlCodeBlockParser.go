package code

import "regexp"

var (
	GraphqlBlockTypeType         CodeBlockType = CodeBlockType{"type", true, regexp.MustCompile(`(?s)^\s*((extend\s+)?\s*type\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 3, 6, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeTypeField    CodeBlockType = CodeBlockType{"type_field", false, regexp.MustCompile(`(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*@[^\n]*)*)\s*$`), 1, 2, 4, []string{"type"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeTypeFieldArg CodeBlockType = CodeBlockType{"type_field_arg", false, regexp.MustCompile(`(?s)\s*((\w+):\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)?)\s*`), 1, 2, -1, []string{"type_field"}, "", "", 0, true, nil}
	GraphqlBlockExplain          CodeBlockType = CodeBlockType{"explain", true, regexp.MustCompile(`(?s)\s*\"\"\"\s*([^\n]*\n)\s*(([^\n]*\s*)*)\s*\"\"\"\s*$`), 0, 1, -1, nil, "", "", 0, false, nil}
	GraphqlBlockExplain2         CodeBlockType = CodeBlockType{"explain2", true, regexp.MustCompile(`(?s)\s*\"([^\n]*)\"\s*$`), 0, 1, -1, nil, "", "", 0, false, nil}
	GraphqlBlockTypeInput        CodeBlockType = CodeBlockType{"input", true, regexp.MustCompile(`(?s)^\s*(extend\s+)?\s*input\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\})\s*$`), 0, 2, 5, nil, "\n", "{}", 4, false, nil}
	GraphqlBlockTypeInputField   CodeBlockType = CodeBlockType{"input_field", true, regexp.MustCompile(`(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)*)\s*$`), 1, 2, 4, []string{"input"}, "\n|,", "()", 3, true, nil}
	GraphqlBlockTypeEnum         CodeBlockType = CodeBlockType{"enum", true, regexp.MustCompile(`(?s)^\s*((extend\s+)?\s*enum\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 3, 6, nil, "\n", "{}", 5, false, nil}
	GraphqlBlockTypeEnumField    CodeBlockType = CodeBlockType{"enum_field", true, regexp.MustCompile(`(?s)^\s*((\w+)(\s*@[^\n]*)*)\s*$`), 1, 2, -1, []string{"enum"}, "\n|,", "()", 1, true, nil}
)

func NewGraphqlCodeBlockParser() *CodeBlockParser {
	return &CodeBlockParser{
		Types: []CodeBlockType{
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
		PairKeys:          []string{"{}", "[]", "()", `""" """`, `""`},
		OriginText:        []string{`""" """`},
		LineCommentKey:    "#",
		PendingLinePrefix: "@",
	}
}
