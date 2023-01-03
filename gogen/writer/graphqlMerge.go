package writer

import (
	"os"

	"github.com/liasece/go-mate/code"
)

func MergeGraphQLFromFile(protoFile string, newContent string) error {
	return mergeGraphQLFromFile(protoFile, newContent)
}

func mergeGraphQLFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := os.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeGraphQL(originFileContent, newContent)
	if toContent != originFileContent {
		// write to file
		err := os.WriteFile(protoFile, []byte(toContent), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func mergeGraphQL(originContent string, newContent string) string {
	c := code.NewGraphqlCodeBlockParser()
	res := c.Parse(originContent).Merge(c.Parse(newContent))
	return res.OriginString
}
