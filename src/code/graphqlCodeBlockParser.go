package code

var (
	GraphqlBlockTypeType         CodeBlockType = CodeBlockType{"type", true, `(?s)^\s*((extend\s+)?\s*type\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\}))\s*$`, 1, 3, 6, nil, "\n", "{}", 5, false}
	GraphqlBlockTypeTypeField    CodeBlockType = CodeBlockType{"type_field", false, `(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*@[^\n]*)*)\s*$`, 1, 2, 4, []string{"type"}, "\n|,", "()", 3, true}
	GraphqlBlockTypeTypeFieldArg CodeBlockType = CodeBlockType{"type_field_arg", false, `(?s)\s*((\w+):\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)?)\s*`, 1, 2, -1, []string{"type_field"}, "", "", 0, true}
	GraphqlBlockExplain          CodeBlockType = CodeBlockType{"explain", true, `(?s)\s*\"\"\"\s*([^\n]*\n)\s*(([^\n]*\s*)*)\s*\"\"\"\s*$`, 0, 1, -1, nil, "", "", 0, false}
	GraphqlBlockExplain2         CodeBlockType = CodeBlockType{"explain2", true, `(?s)\s*\"([^\n]*)\"\s*$`, 0, 1, -1, nil, "", "", 0, false}
	GraphqlBlockTypeInput        CodeBlockType = CodeBlockType{"input", true, `(?s)^\s*(extend\s+)?\s*input\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{\s*(.*?)\s*\})\s*$`, 0, 2, 5, nil, "\n", "{}", 4, false}
	GraphqlBlockTypeInputField   CodeBlockType = CodeBlockType{"input_field", true, `(?s)^\s*((\w+)(\(([\S\s]+)\))?:\s*([\w!\[\]]+)(\s*=\s*[\w]*)?(\s*@[^\n]*)*)\s*$`, 1, 2, 4, []string{"input"}, "\n|,", "()", 3, true}
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
		},
		PairKeys:          []string{"{}", "[]", "()", `""" """`, "\"\""},
		LineCommentKey:    "#",
		PendingLinePrefix: "@",
	}
}
