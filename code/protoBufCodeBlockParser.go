package code

import "regexp"

var (
	ProtoBlockTypeImport       = BlockType{"import", regexp.MustCompile(`(?s)^\s*(import\s*(.+?);?\n?)\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, nil, "", "", 1, false, nil}
	ProtoBlockTypePackage      = BlockType{"package", regexp.MustCompile(`(?s)^\s*(package\s*(.+?);?\n?)\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, nil, "", "", 1, false, nil}
	ProtoBlockTypeSyntax       = BlockType{"syntax", regexp.MustCompile(`(?s)^\s*(syntax\s*=\s*(.+?);?\n?)\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, nil, "", "", 1, false, nil}
	ProtoBlockTypeService      = BlockType{"service", regexp.MustCompile(`(?s)^\s*(service\s+(\w+)\s*(\{\s?(.*?\n?)\s*\}\n?))\s*$`), 1, 2, []int{4}, [][]string{nil}, []bool{true}, nil, "\n", "{}", 3, false, []string{";", "}"}}
	ProtoBlockTypeRPC          = BlockType{"rpc", regexp.MustCompile(`(?s)^\s*(rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{\n?(\s*.*?\n?)\s*\})?\s*?;?\n?(\s*//.*)?)\s*$`), 1, 2, []int{4}, [][]string{nil}, []bool{false}, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeMessage      = BlockType{"message", regexp.MustCompile(`(?s)^\s*(message\s+(\w+)\s*(\{\n?(.*?\n?)\s*\})\s*?;?\n?)\s*$`), 1, 2, []int{4}, [][]string{nil}, []bool{true}, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeMessageField = BlockType{"message_field", regexp.MustCompile(`(?s)\s*((optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?\n?(\s*//.*)?)\s*?`), 1, 5, []int{8}, [][]string{nil}, []bool{true}, nil, ",|\n", "[]", 7, false, nil}
	ProtoBlockTypeOption       = BlockType{"option", regexp.MustCompile(`(?s)^\s*(option\s+(.+?)\s*=\s*(\{?\s*(.*?)\s*\}?)?\s*?;?\n?)\s*$`), 1, 2, []int{4}, [][]string{nil}, []bool{true}, nil, ",", "{}", 3, false, nil}
	ProtoBlockTypeOptionItem   = BlockType{"option_item", regexp.MustCompile(`(?s)(\s*(\w+)\s*[:=]\s*([^,]+))\s*`), 1, 2, []int{}, [][]string{}, []bool{}, []string{"option", "message_field"}, "", "", 0, false, nil}
	ProtoBlockTypeOptionItem2  = BlockType{"option_item2", regexp.MustCompile(`(?s)\s*(\"([^\n]*)\")\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, []string{"option", "message_field"}, "", "", 0, false, nil}
	ProtoBlockTypeEnum         = BlockType{"enum", regexp.MustCompile(`(?s)^\s*(enum\s*(.+?)\s*(\{\s*(.*)\s*\})\s*?;?\n?)\s*$`), 1, 2, []int{4}, [][]string{nil}, []bool{true}, nil, "\n", "{}", 3, false, nil}
	ProtoBlockTypeEnumItem     = BlockType{"enum_item", regexp.MustCompile(`(?s)\s*((\w+)\s*[:=]\s*(\S+))\s*`), 1, 2, []int{}, [][]string{}, []bool{}, []string{"enum"}, "", "", 1, false, nil}
	ProtoBlockTypeReserved     = BlockType{"reserved", regexp.MustCompile(`(?s)^\s*(reserved\s*(.+?);?\n?)\s*$`), 1, 2, []int{}, [][]string{}, []bool{}, nil, "", "", 1, false, nil}
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
		PairKeys:           []string{"{}", "[]", "()"},
		LineCommentKey:     "//",
		OriginText:         nil,
		PendingLinePrefix:  "",
		HeadViscousPairKey: nil,
	}
}
