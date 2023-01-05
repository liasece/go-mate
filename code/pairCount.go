package code

import (
	"strings"
)

type PairCount struct {
	Count      map[string]int // like "()":1 "{}":1 "[]":1
	KeyWord    []string       // like () {} [], in one string, first rune is left, second rune is right
	OriginText []string
}

func pairKeySplit(key string) (head string, tail string) {
	if len(key) == 2 {
		head = key[:1]
		tail = key[1:]
	} else {
		list := strings.Split(key, " ")
		if len(list) == 2 {
			head = list[0]
			tail = list[1]
		} else {
			panic("CodePairCount.Add: invalid keyword: " + key)
		}
	}
	return
}

func (c *PairCount) Add(line string) (effectKey string) {
	inOrigin := ""
	for k, v := range c.Count {
		if v > 0 {
			for _, originWarp := range c.OriginText {
				if k == originWarp {
					inOrigin = originWarp
					break
				}
			}
		}
	}
	for _, key := range c.KeyWord {
		if inOrigin != "" {
			if key != inOrigin {
				continue
			}
		}
		head, tail := pairKeySplit(key)
		if head == tail {
			if count := strings.Count(line, head); count > 0 {
				if c.Count[key] > 0 {
					c.Count[key] -= count % 2
				} else {
					c.Count[key] += count % 2
				}
				effectKey = key
			}
		} else {
			headCount := strings.Count(line, head)
			if headCount > 0 {
				c.Count[key] += headCount
				effectKey = key
			}
			tailCount := strings.Count(line, tail)
			if tailCount > 0 {
				c.Count[key] -= tailCount
				effectKey = key
			}
		}

		// no repeated count
		line = strings.ReplaceAll(line, head, "")
		line = strings.ReplaceAll(line, tail, "")
	}
	return effectKey
}

func (c *PairCount) IsZero() bool {
	for _, v := range c.KeyWord {
		if c.Count[v] != 0 {
			return false
		}
	}
	return true
}
