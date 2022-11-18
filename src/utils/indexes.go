package utils

import (
	"context"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/liasece/gocoder"
	"gorm.io/gorm/schema"
)

type Index struct {
	Name    string
	Class   string // UNIQUE | FULLTEXT | SPATIAL
	Type    string // btree, hash, gist, spgist, gin, and brin
	Where   string
	Comment string
	Option  string // WITH PARSER parser_name
	Fields  []IndexOption
}

type IndexOption struct {
	gocoder.Field
	Expression string
	Sort       string // DESC, ASC
	Collate    string
	Length     int
	Priority   int
}

func (o *IndexOption) GetTag(tag string) string {
	return reflect.StructTag(o.Field.GetTag()).Get(tag)
}

func (o *IndexOption) Name() string {
	return o.Field.GetName()
}

func parseFieldIndexes(field gocoder.Field) (indexes []*Index) {
	tagStr := reflect.StructTag(field.GetTag()).Get("gorm")
	// log.L(context.Background()).Error("parseFieldIndexes", log.Any("tagStr", tagStr), log.Any("field.GetTag()", field.GetTag()))
	for _, value := range strings.Split(tagStr, ";") {
		if value != "" {
			v := strings.Split(value, ":")
			k := strings.TrimSpace(strings.ToUpper(v[0]))
			if k == "INDEX" || k == "UNIQUEINDEX" {
				var (
					name      string
					tag       = strings.Join(v[1:], ":")
					idx       = strings.Index(tag, ",")
					settings  = schema.ParseTagSetting(tag, ",")
					length, _ = strconv.Atoi(settings["LENGTH"])
				)

				if idx == -1 {
					idx = len(tag)
				}

				if idx != -1 {
					name = tag[0:idx]
				}

				if (k == "UNIQUEINDEX") || settings["UNIQUE"] != "" {
					settings["CLASS"] = "UNIQUE"
				}

				priority, err := strconv.Atoi(settings["PRIORITY"])
				if err != nil {
					priority = 10
				}

				indexes = append(indexes, &Index{
					Name:    name,
					Class:   settings["CLASS"],
					Type:    settings["TYPE"],
					Where:   settings["WHERE"],
					Comment: settings["COMMENT"],
					Option:  settings["OPTION"],
					Fields: []IndexOption{{
						Field:      field,
						Expression: settings["EXPRESSION"],
						Sort:       settings["SORT"],
						Collate:    settings["COLLATE"],
						Length:     length,
						Priority:   priority,
					}},
				})
			}
		}
	}

	return
}

// get index config from gorm tag
func GormIndexes(ctx context.Context, t gocoder.Type) []*Index {
	var indexes = []*Index{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		for _, index := range parseFieldIndexes(field) {
			idx := &Index{}
			find := false
			for _, v := range indexes {
				if v.Name == index.Name {
					idx = v
					find = true
					break
				}
			}
			if !find {
				indexes = append(indexes, idx)
			}
			idx.Name = index.Name
			if idx.Class == "" {
				idx.Class = index.Class
			}
			if idx.Type == "" {
				idx.Type = index.Type
			}
			if idx.Where == "" {
				idx.Where = index.Where
			}
			if idx.Comment == "" {
				idx.Comment = index.Comment
			}
			if idx.Option == "" {
				idx.Option = index.Option
			}

			idx.Fields = append(idx.Fields, index.Fields...)
			sort.Slice(idx.Fields, func(i, j int) bool {
				return idx.Fields[i].Priority < idx.Fields[j].Priority
			})
		}
	}

	// for _, idx := range indexes {
	// 	log.L(ctx).Error("index", log.Any("idx", idx))
	// }
	// log.L(ctx).Panic("GormIndexes finish")
	return indexes
}
