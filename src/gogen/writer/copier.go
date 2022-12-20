package writer

import (
	"bufio"
	"io/ioutil"
	"regexp"
	"strings"
)

func FillCopierLine(copierFile string, names [][2]string) error {
	originFileContent := ""
	{
		// read from file
		content, err := ioutil.ReadFile(copierFile)
		if err == nil {
			originFileContent = string(content)
		}
	}
	toContent := buildCopierContent(originFileContent, names)
	if toContent != originFileContent {
		// write to file
		err := ioutil.WriteFile(copierFile, []byte(toContent), 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCopierFromStr(originContent string, nameSrc string, nameDest string) string {
	scanner := bufio.NewScanner(strings.NewReader(originContent))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	out := ""
	for _, l := range lines {
		nameReg := regexp.MustCompile(`.*?` + regexp.QuoteMeta(nameSrc) + `\).*?` + regexp.QuoteMeta(nameDest) + `\).*`)
		parts := nameReg.FindStringSubmatch(l)
		if len(parts) != 0 {
			return l
		}
	}
	return out
}

func cutLastCopierFromStr(originContent string) (string, string) {
	scanner := bufio.NewScanner(strings.NewReader(originContent))
	lines := make([]string, 0)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	head := ""
	tail := ""
	lastIn := false
	isTail := false
	for _, l := range lines {
		nameReg := regexp.MustCompile(`\s+g.Add\(\(`)
		parts := nameReg.FindStringSubmatch(l)
		if len(parts) != 0 {
			lastIn = true
		} else {
			if lastIn {
				isTail = true
			}
		}
		if isTail {
			tail += l + "\n"
		} else {
			head += l + "\n"
		}
	}
	return head, tail
}

func buildCopierContent(originContent string, names [][2]string) string {
	res := originContent
	for _, name := range names {
		matchOrigin := getCopierFromStr(res, name[0], name[1])
		addFsStr := `	g.Add((*` + name[0] + `)(nil), (*` + name[1] + `)(nil))`
		if matchOrigin != "" {
			res = strings.Replace(res, matchOrigin, addFsStr, 1)
		} else {
			head, tail := cutLastCopierFromStr(res)
			res = head + addFsStr + "\n" + tail
		}
	}
	return res
}
