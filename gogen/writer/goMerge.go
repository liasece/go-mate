package writer

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/tools/imports"
)

func splitGoBlock(content string) (blocks []string, body []string) {
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

func getGoBlockHead(blockContent string) string {
	nameReg := regexp.MustCompile(`(.*?[\{\(])\n`)
	parts := nameReg.FindStringSubmatch(blockContent)
	if len(parts) > 0 {
		return parts[1]
	}
	return blockContent
}

func getGoBlockByHead(blocks []string, head string) (string, int) {
	emptyReg := regexp.MustCompile(`\s+`)
	regRule := `\s*` + emptyReg.ReplaceAllString(regexp.QuoteMeta(head), `\s+`)
	// log.Info("getGoBlockByHead begin: "+regRule, log.Any("head", head), log.Any("regRule", regRule))
	headReg := regexp.MustCompile(regRule)
	for i, b := range blocks {
		parts := headReg.FindStringSubmatch(b)
		if len(parts) != 0 {
			return b, i
		}
	}
	return "", -1
}

func MergeGoFromFile(protoFile string, newContent string) error {
	return mergeGoFromFile(protoFile, newContent)
}

func mergeGoFromFile(protoFile string, newContent string) error {
	originFileContent := ""
	{
		// read from file
		content, err := os.ReadFile(protoFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := mergeGo(originFileContent, newContent)
	if toContent != originFileContent {
		bytes, err := imports.Process(protoFile, []byte(toContent), &imports.Options{
			FormatOnly: true,
			Comments:   true,
			TabIndent:  true,
			TabWidth:   8,
			Fragment:   false,
			AllErrors:  false,
		})
		if err != nil {
			return err
		}
		err = os.WriteFile(protoFile, bytes, 0600)
		if err != nil {
			return errors.Wrapf(err, "failed to write %s", protoFile)
		}
		// log.Error("mergeGoFromFile finish, changed", log.Any("protoFile", protoFile))
		// if !strings.HasPrefix(toContent, "//go:build wireinject") {
		// 	log.Error(toContent)
		// 	os.Exit(1)
		// }
	}
	return nil
}

func mergeGo(originContent string, newContent string) string {
	newBlocks, newBody := splitGoBlock(newContent)
	originBlocks, originBody := splitGoBlock(originContent)
	// log.Error("mergeGo", log.Any("newContent", newContent), log.Any("originContent", originContent), log.Any("newBlocks", newBlocks), log.Any("newBody", newBody), log.Any("originBlocks", originBlocks), log.Any("originBody", originBody))
	res := originContent
	for i, newBlock := range newBlocks {
		if newBlock == "\n" {
			continue
		}
		newBlockHead := getGoBlockHead(newBlock)
		origin, index := getGoBlockByHead(originBlocks, newBlockHead)
		// log.Info("mergeGo", log.Any("b", newBlock), log.Any("newHead", newBlockHead), log.Any("origin", origin), log.Any("index", index))
		if origin == "" {
			// add
			res += newBlock
		} else {
			// replace
			if strings.Count(newBlock, "\n") > 1 {
				oldContent := originBody[index]
				newContent := mergeGo(originBody[index], newBody[i])
				res = strings.Replace(res, oldContent, newContent, 1)
			} else {
				res = strings.Replace(res, origin, newBlock, 1)
			}
		}
	}
	// log.Error("mergeGo finish: " + res)
	// if strings.Contains(res, "Collection") {
	// 	os.Exit(1)
	// }
	return res
}
