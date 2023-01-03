package code

import "regexp"

var (
	ProtoBlockTypeImport       = BlockType{"import", false, regexp.MustCompile(`(?s)^\s*(import\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false, nil}
	ProtoBlockTypePackage      = BlockType{"package", false, regexp.MustCompile(`(?s)^\s*(package\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false, nil}
	ProtoBlockTypeSyntax       = BlockType{"syntax", false, regexp.MustCompile(`(?s)^\s*(syntax\s*=\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false, nil}
	ProtoBlockTypeService      = BlockType{"service", true, regexp.MustCompile(`(?s)^\s*(service\s+(\w+)\s*(\{\s*(.*?)\s*\}))\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false, []string{";", "}"}}
	ProtoBlockTypeRPC          = BlockType{"rpc", false, regexp.MustCompile(`(?s)^\s*(rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{\s*(.*?)\s*\})?\s*?;?(\s*//.*)?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeMessage      = BlockType{"message", true, regexp.MustCompile(`(?s)^\s*(message\s+(\w+)\s*(\{\s*(.*?)\s*\})\s*?;?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeMessageField = BlockType{"message_field", true, regexp.MustCompile(`(?s)\s*((optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?(\s*//.*)?)\s*?`), 1, 5, 8, nil, ",|\n", "[]", 7, false, nil}
	ProtoBlockTypeOption       = BlockType{"option", true, regexp.MustCompile(`(?s)^\s*(option\s+(.+?)\s*=\s*(\{?\s*(.*?)\s*\}?)?\s*?;?)\s*$`), 1, 2, 4, nil, ",", "{}", 3, false, nil}
	ProtoBlockTypeOptionItem   = BlockType{"option_item", false, regexp.MustCompile(`(?s)\s*((\w+)\s*[:=]\s*([^,]+))\s*[, }\n]?\s*`), 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false, nil}
	ProtoBlockTypeOptionItem2  = BlockType{"option_item2", false, regexp.MustCompile(`(?s)\s*(\"([^\n]*)\")\s*$`), 1, 2, -1, []string{"option", "message_field"}, "", "", 0, false, nil}
	ProtoBlockTypeEnum         = BlockType{"enum", true, regexp.MustCompile(`(?s)^\s*(enum\s*(.+?)\s*(\{\s*(.*)\s*\})\s*?;?)\s*$`), 1, 2, 4, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeEnumItem     = BlockType{"enum_item", false, regexp.MustCompile(`(?s)\s*((\w+)\s*[:=]\s*(\S+))\s*`), 1, 2, -1, []string{"enum"}, "", "", 1, false, nil}
	ProtoBlockTypeReserved     = BlockType{"reserved", true, regexp.MustCompile(`(?s)^\s*(reserved\s*(.+?);?)\s*$`), 1, 2, -1, nil, "", "", 1, false, nil}
)

func NewProtoBufCodeBlockParser() *BlockParser {
	return &BlockParser{
		Types: []BlockType{
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
		PairKeys:          []string{"{}", "[]", "()"},
		LineCommentKey:    "//",
		OriginText:        nil,
		PendingLinePrefix: "",
	}
}
