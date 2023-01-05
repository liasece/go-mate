package writer

import (
	"os"

	"github.com/liasece/go-mate/code"
)

func MergeProtoFromFile(protoFile string, newContent string) error {
	return mergeProtoFromFile(protoFile, newContent)
}

func mergeProtoFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := os.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeProto(originFileContent, newContent)
	if toContent != originFileContent {
		// write to file
		err := os.WriteFile(protoFile, []byte(toContent), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func MergeProto(originContent string, newContent string) string {
	return mergeProto(originContent, newContent)
}

func mergeProto(originContent string, newContent string) string {
	c := code.NewProtoBufCodeBlockParser()
	res := c.Parse(originContent).Merge(0, c.Parse(newContent))
	toContent := res.OriginString
	{
		// add end line
		if toContent != "" && toContent[len(toContent)-1] != '\n' {
			toContent += "\n"
		}
	}
	return toContent
}
