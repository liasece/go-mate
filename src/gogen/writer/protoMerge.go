package writer

import (
	"io/ioutil"

	"github.com/liasece/go-mate/src/code"
)

func MergeProtoFromFile(protoFile string, newContent string) error {
	return mergeProtoFromFile(protoFile, newContent)
}

func mergeProtoFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := ioutil.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeProto(originFileContent, newContent)
	if toContent != originFileContent {
		// write to file
		err := ioutil.WriteFile(protoFile, []byte(toContent), 0644)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeProto(originContent string, newContent string) string {
	c := code.NewProtoBufCodeBlockParser()
	res := c.Parse(originContent).Merge(c.Parse(newContent))
	return res.OriginString
}
