package writer

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/liasece/go-mate/src/code"
)

func splitGraphQLBlock(content string) (blocks []string, body []string) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	tmpOut := ""
	tmpBody := ""
	in := false
	res := make([]string, 0)
	resBody := make([]string, 0)
	for i := 0; i < len(lines); i++ {
		l := lines[i]
		isHead := false
		if !in {
			nameReg := regexp.MustCompile(`^.*?[\({]$`)
			parts := nameReg.FindStringSubmatch(l)
			if len(parts) != 0 {
				in = true
				isHead = true
			} else {
				// get go line block
				block := l + "\n"
				res = append(res, block)
				resBody = append(resBody, block)
			}
		}
		if in {
			tmpOut += l + "\n"

			nameReg := regexp.MustCompile(`^[)}]`)
			parts := nameReg.FindStringSubmatch(l)
			if len(parts) != 0 {
				in = false
				res = append(res, tmpOut)
				resBody = append(resBody, tmpBody)
				tmpOut = ""
				tmpBody = ""
			} else if !isHead {
				tmpBody += l + "\n"
			}
		}
	}
	if tmpOut != "" {
		res = append(res, tmpOut)
	}
	if tmpBody != "" {
		resBody = append(resBody, tmpBody)
	}
	return res, resBody
}

func getGraphQLBlockHead(blockContent string) string {
	nameReg := regexp.MustCompile(`(.*?[\{\(])\n`)
	parts := nameReg.FindStringSubmatch(blockContent)
	if len(parts) > 0 {
		return parts[1]
	}
	return blockContent
}

func getGraphQLBlockByHead(blocks []string, head string) (string, int) {
	regRule := `\s*` + strings.ReplaceAll(regexp.QuoteMeta(head), " ", `\s+`)
	// log.Info("getGraphQLBlockByHead begin: "+regRule, log.Any("head", head), log.Any("regRule", regRule))
	headReg := regexp.MustCompile(regRule)
	for i, b := range blocks {
		parts := headReg.FindStringSubmatch(b)
		if len(parts) != 0 {
			return b, i
		}
	}
	return "", -1
}

func MergeGraphQLFromFile(protoFile string, newContent string) error {
	return mergeGraphQLFromFile(protoFile, newContent)
}

func mergeGraphQLFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := ioutil.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeGraphQL(originFileContent, newContent)
	if toContent != originFileContent {
		// write to file
		err := ioutil.WriteFile(protoFile, []byte(toContent), 0600)
		if err != nil {
			return err
		}
		// log.Error("mergeGraphQLFromFile finish, changed", log.Any("protoFile", protoFile))
		// if !strings.HasPrefix(toContent, "//go:build wireinject") {
		// 	log.Error(toContent)
		// 	os.Exit(1)
		// }
	}
	return nil
}

func mergeGraphQL(originContent string, newContent string) string {
	c := code.NewGraphqlCodeBlockParser()
	res := c.Parse(originContent).Merge(c.Parse(newContent))
	return res.OriginString
}
