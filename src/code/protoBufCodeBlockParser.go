package code

import "regexp"

var (
	ProtoBlockTypeImport       CodeBlockType = CodeBlockType{"import", false, regexp.MustCompile(`(?s)^\s*(import\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypePackage      CodeBlockType = CodeBlockType{"package", false, regexp.MustCompile(`(?s)^\s*(package\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypeSyntax       CodeBlockType = CodeBlockType{"syntax", false, regexp.MustCompile(`(?s)^\s*(syntax\s*=\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false}
	ProtoBlockTypeService      CodeBlockType = CodeBlockType{"service", true, regexp.MustCompile(`(?s)^\s*(service\s+(\w+)\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeRPC          CodeBlockType = CodeBlockType{"rpc", false, regexp.MustCompile(`(?s)^\s*(rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{\s*(.*?)\s*\})?\s*?;?(\s*//.*)?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeMessage      CodeBlockType = CodeBlockType{"message", true, regexp.MustCompile(`(?s)^\s*(message\s+(\w+)\s*(\{\s*(.*?)\s*\})\s*?;?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeMessageField CodeBlockType = CodeBlockType{"message_field", true, regexp.MustCompile(`(?s)\s*((optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?(\s*//.*)?)\s*?`), 1, 5, 8, nil, ",|\n", "[]", 7, false}
	ProtoBlockTypeOption       CodeBlockType = CodeBlockType{"option", true, regexp.MustCompile(`(?s)^\s*(option\s+(.+?)\s*=\s*(\{?\s*(.*?)\s*\}?)?\s*?;?)\s*$`), 1, 2, 4, nil, ",", "{}", 3, false}
	ProtoBlockTypeOptionItem   CodeBlockType = CodeBlockType{"option_item", false, regexp.MustCompile(`(?s)\s*((\w+)\s*[:=]\s*([^,]+))\s*[, }\n]?\s*`), 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false}
	ProtoBlockTypeOptionItem2  CodeBlockType = CodeBlockType{"option_item2", false, regexp.MustCompile(`(?s)\s*(\"([^\n]*)\")\s*$`), 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false}
	ProtoBlockTypeEnum         CodeBlockType = CodeBlockType{"enum", true, regexp.MustCompile(`(?s)^\s*(enum\s*(.+?)\s*(\{\s*(.*)\s*\})\s*?;?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false}
	ProtoBlockTypeEnumItem     CodeBlockType = CodeBlockType{"enum_item", false, regexp.MustCompile(`(?s)\s*((\w+)\s*[:=]\s*(\S+))\s*`), 1, 2, -1, []string{"enum"}, "", "", 1, false}
	ProtoBlockTypeReserved     CodeBlockType = CodeBlockType{"reserved", true, regexp.MustCompile(`(?s)^\s*(reserved\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false}
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
