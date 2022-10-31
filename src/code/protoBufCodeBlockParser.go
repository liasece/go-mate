package code

var (
	ProtoBlockTypeNone         CodeBlockType = CodeBlockType{"", true, "", 0, 0, 0, nil}
	ProtoBlockTypeImport       CodeBlockType = CodeBlockType{"import", false, `(?s)^\s*import\s*(.+?);?\s*$`, 0, 1, -1, nil}
	ProtoBlockTypePackage      CodeBlockType = CodeBlockType{"package", false, `(?s)^\s*package\s*(.+?);?\s*$`, 0, 1, -1, nil}
	ProtoBlockTypeSyntax       CodeBlockType = CodeBlockType{"syntax", false, `(?s)^\s*syntax\s*=\s*(.+?);?\s*$`, 0, 1, -1, nil}
	ProtoBlockTypeService      CodeBlockType = CodeBlockType{"service", true, `(?s)^\s*service\s+(\w+)\s*\{(.*?)\}\s*$`, 0, 1, 2, nil}
	ProtoBlockTypeRPC          CodeBlockType = CodeBlockType{"rpc", false, `(?s)^\s*rpc\s+(\w+)\s*\(.*?\)\s*returns\s*\(.*?\)\s*(\{(.*)\})?\s*;?(\s*//.*)?\s*$`, 0, 1, 3, nil}
	ProtoBlockTypeMessage      CodeBlockType = CodeBlockType{"message", true, `(?s)^\s*message\s+(\w+)\s*\{(.*)\}\s*;?\s*$`, 0, 1, 2, nil}
	ProtoBlockTypeMessageField CodeBlockType = CodeBlockType{"message_field", true, `(?s)^\s*(optional)?(repeated)?\s*([a-zA-Z0-9.<>, ]+)?\s+(\w+)\s*=\s*(\d+)\s*(\[(.*?)\])?;?(\s*//.*)?\s*?$`, 0, 4, 7, nil}
	ProtoBlockTypeOption       CodeBlockType = CodeBlockType{"option", true, `(?s)^\s*option\s+(.+?)\s*=\s*\{?(.*?)\s*\}?\s*;?\s*$`, 0, 1, 2, nil}
	ProtoBlockTypeOptionItem   CodeBlockType = CodeBlockType{"option_item", false, `(?s)\s*((\w+)\s*[:=]\s*([^,]+))\s*[, }\n]?\s*`, 1, 2, -1, []string{"option", "message_field"}}
	ProtoBlockTypeEnum         CodeBlockType = CodeBlockType{"enum", true, `(?s)^\s*enum\s*(.+?)\s*\{(.*)\}\s*;?\s*$`, 0, 1, 2, nil}
	ProtoBlockTypeEnumItem     CodeBlockType = CodeBlockType{"enum_item", false, `(?s)\s*(\w+)\s*[:=]\s*(\S+)\s*`, 0, 1, -1, []string{"enum"}}
	ProtoBlockTypeReserved     CodeBlockType = CodeBlockType{"reserved", true, `(?s)^\s*reserved\s*(.+?);?\s*$`, 0, 1, -1, nil}
)

func NewProtoBufCodeBlockParser() *CodeBlockParser {
	return &CodeBlockParser{
		Types: []CodeBlockType{
			ProtoBlockTypeNone,
			ProtoBlockTypeImport,
			ProtoBlockTypePackage,
			ProtoBlockTypeSyntax,
			ProtoBlockTypeService,
			ProtoBlockTypeRPC,
			ProtoBlockTypeMessage,
			ProtoBlockTypeMessageField,
			ProtoBlockTypeOption,
			ProtoBlockTypeOptionItem,
			ProtoBlockTypeEnum,
			ProtoBlockTypeEnumItem,
			ProtoBlockTypeReserved,
		},
		PairKeys:       []string{"{}", "[]", "()"},
		LineCommentKey: "//",
	}
}
