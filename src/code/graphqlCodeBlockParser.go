package code

var (
	GraphqlBlockTypeNone         CodeBlockType = CodeBlockType{"", true, "", 0, 0, 0, nil, "", "", 0}
	GraphqlBlockTypeType         CodeBlockType = CodeBlockType{"type", true, `(?s)^\s*(extend\s+)?\s*type\s+(\w+)\s*(implements\s+\w+\s*)?\s*(\{(.*?)\})\s*$`, 0, 2, 5, nil, "", "{}", 4}
	GraphqlBlockTypeTypeField    CodeBlockType = CodeBlockType{"type_field", true, `(?s)^\s*(\w+)(\(([\w!\[\] ,:]+)\))?:\s*([\w!\[\]]+)(\s*@[^\n]*)?\s*$`, 0, 1, 3, nil, ",", "()", 2}
	GraphqlBlockTypeTypeFieldArg CodeBlockType = CodeBlockType{"type_field_arg", true, `(?s)\s*(\w+):\s*([\w!\[\]]+)(\s*@[^\n]*)?\s*`, 0, 1, -1, []string{"type_field"}, "", "", 0}
)

func NewGraphqlBufCodeBlockParser() *CodeBlockParser {
	return &CodeBlockParser{
		Types: []CodeBlockType{
			GraphqlBlockTypeNone,
			GraphqlBlockTypeType,
			GraphqlBlockTypeTypeField,
			GraphqlBlockTypeTypeFieldArg,
		},
		PairKeys:       []string{"{}", "[]", "()"},
		LineCommentKey: "#",
	}
}
