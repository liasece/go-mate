package writer

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strings"
)

func splitProtoBlock(content string) (blocks []string, body []string) {
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
			nameReg := regexp.MustCompile(`^\s*[a-zA-Z_ \t]*?\s*{`)
			parts := nameReg.FindStringSubmatch(l)
			if len(parts) != 0 {
				in = true
				isHead = true
			} else {
				nameReg := regexp.MustCompile(`^\s*rpc\s+.*?\s*\(`)
				parts := nameReg.FindStringSubmatch(l)
				if len(parts) != 0 {
					// get rpc line block
					block := l + "\n"
					{
						// if mul line
						if parts := regexp.MustCompile(`;\s*$`).FindStringSubmatch(l); len(parts) == 0 {
							for ; i < len(lines); i++ {
								l := lines[i]
								block += l + "\n"
								if parts := regexp.MustCompile(`}\s*;?\s*$`).FindStringSubmatch(l); len(parts) != 0 {
									break
								}
							}
						}
					}
					res = append(res, block)
					resBody = append(resBody, block)
				}
			}
		}
		if in {
			tmpOut += l + "\n"

			nameReg := regexp.MustCompile(`^}`)
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

func getProtoBlockHead(blockContent string) string {
	nameReg := regexp.MustCompile(`\s*(\w+)\s+(\w+)\s*[\{\(]`)
	parts := nameReg.FindStringSubmatch(blockContent)
	if len(parts) == 3 {
		return parts[1] + " " + parts[2]
	}
	return ""
}

func getProtoBlockByHead(blocks []string, head string) (string, int) {
	headReg := regexp.MustCompile(`\s*` + strings.ReplaceAll(head, " ", `\s+`) + `\s*[{\()]`)
	for i, b := range blocks {
		parts := headReg.FindStringSubmatch(b)
		if len(parts) != 0 {
			return b, i
		}
	}
	return "", -1
}

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
	newBlocks, newBody := splitProtoBlock(newContent)
	originBlocks, originBody := splitProtoBlock(originContent)
	// log.Info("mergeProto", log.Any("newContent", newContent), log.Any("originContent", originContent), log.Any("newBlocks", newBlocks), log.Any("newBody", newBody), log.Any("originBlocks", originBlocks), log.Any("originBody", originBody))
	res := originContent
	for i, b := range newBlocks {
		newHead := getProtoBlockHead(b)
		// log.Info("mergeProto", log.Any("newHead", newHead))
		if origin, index := getProtoBlockByHead(originBlocks, newHead); origin == "" {
			// add
			res = res + b
		} else {
			// replace
			if strings.HasPrefix(newHead, "service") {
				res = strings.Replace(res, originBody[index], mergeProto(originBody[index], newBody[i]), 1)
			} else {
				res = strings.Replace(res, origin, b, 1)
			}
		}
	}
	return res
}
