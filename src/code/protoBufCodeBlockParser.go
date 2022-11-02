package code

var (
	ProtoBlockTypeImport       CodeBlockType = CodeBlockType{"import", false, `(?s)^\s*(import\s*(.+?);?)\s*$`, 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypePackage      CodeBlockType = CodeBlockType{"package", false, `(?s)^\s*(package\s*(.+?);?)\s*$`, 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypeSyntax       CodeBlockType = CodeBlockType{"syntax", false, `(?s)^\s*(syntax\s*=\s*(.+?);?)\s*$`, 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypeService      CodeBlockType = CodeBlockType{"service", true, `(?s)^\s*(service\s+(\w+)\s*(\{\s*(.*?)\s*\}))\s*$`, 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeRPC          CodeBlockType = CodeBlockType{"rpc", false, `(?s)^\s*(rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{\s*(.*?)\s*\})?\s*?;?(\s*//.*)?)\s*$`, 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeMessage      CodeBlockType = CodeBlockType{"message", true, `(?s)^\s*(message\s+(\w+)\s*(\{\s*(.*?)\s*\})\s*?;?)\s*$`, 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeMessageField CodeBlockType = CodeBlockType{"message_field", true, `(?s)\s*((optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?(\s*//.*)?)\s*?`, 1, 5, 8, nil, ",|\n", "[]", 7, false}
	ProtoBlockTypeOption       CodeBlockType = CodeBlockType{"option", true, `(?s)^\s*(option\s+(.+?)\s*=\s*(\{?\s*(.*?)\s*\}?)?\s*?;?)\s*$`, 1, 2, 4, nil, ",", "{}", 3, false}
	ProtoBlockTypeOptionItem   CodeBlockType = CodeBlockType{"option_item", false, `(?s)\s*((\w+)\s*[:=]\s*([^,]+))\s*[, }\n]?\s*`, 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false}
	ProtoBlockTypeOptionItem2  CodeBlockType = CodeBlockType{"option_item2", false, `(?s)\s*(\"([^\n]*)\")\s*$`, 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false}
	ProtoBlockTypeEnum         CodeBlockType = CodeBlockType{"enum", true, `(?s)^\s*(enum\s*(.+?)\s*(\{\s*(.*)\s*\})\s*?;?)\s*$`, 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeEnumItem     CodeBlockType = CodeBlockType{"enum_item", false, `(?s)\s*((\w+)\s*[:=]\s*(\S+))\s*`, 1, 2, -1, []string{"enum"}, "", "", 1, false}
	ProtoBlockTypeReserved     CodeBlockType = CodeBlockType{"reserved", true, `(?s)^\s*(reserved\s*(.+?);?)\s*$`, 1, 2, -1, nil, "", "", 1, false}
)

func NewProtoBufCodeBlockParser() *CodeBlockParser {
	return &CodeBlockParser{
		Types: []CodeBlockType{
			ProtoBlockTypeImport,
			ProtoBlockTypePackage,
			ProtoBlockTypeSyntax,
			ProtoBlockTypeService,
			ProtoBlockTypeRPC,
			ProtoBlockTypeMessage,
			ProtoBlockTypeMessageField,
			ProtoBlockTypeOption,
			ProtoBlockTypeOptionItem,
			ProtoBlockTypeOptionItem2,
			ProtoBlockTypeEnum,
			ProtoBlockTypeEnumItem,
			ProtoBlockTypeReserved,
		},
		PairKeys:       []string{"{}", "[]", "()"},
		LineCommentKey: "//",
	}
}
