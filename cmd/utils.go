package cmd

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/liasece/log"
	"github.com/spf13/cobra"
)

func calGoFilePkgName(path string) string {
	if fileInfo, err := os.Stat(path); err == nil {
		if fileInfo.IsDir() {
			fset := token.NewFileSet()
			// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
			fm, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
			if err == nil {
				for k, v := range fm {
					if k != "" {
						if v.Name != "" {
							return v.Name
						}
						return k
					}
				}
			}
		} else {
			fset := token.NewFileSet()
			// 这里取绝对路径，方便打印出来的语法树可以转跳到编辑器
			f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
			if err == nil && f.Name.Name != "" {
				return f.Name.Name
			}
		}
	}
	dir, f := filepath.Split(path)
	if strings.Contains(f, ".") {
		return filepath.Base(dir)
	}
	if f == "" {
		return filepath.Base(dir)
	}
	return f
}

func InitFlag(c *cobra.Command, cfg interface{}) {
	// build arg flag
	v := reflect.ValueOf(cfg)
	t := v.Type().Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fv := v.Elem().Field(i)
		tag := f.Tag.Get("arg")
		tagS := strings.Split(tag, ";")
		tagMap := make(map[string]string)
		for _, item := range tagS {
			item = strings.Trim(item, " ")
			kvs := strings.Split(item, ":")
			if len(kvs) == 1 {
				tagMap[kvs[0]] = ""
			} else if len(kvs) == 2 {
				v := strings.Trim(kvs[1], " ")
				tagMap[kvs[0]] = v
			}
		}
		name := f.Name
		if v, ok := tagMap["name"]; ok && v != "" {
			name = v
		}
		short := ""
		if v, ok := tagMap["short"]; ok && v != "" {
			short = v
		}
		defaultV := ""
		if v, ok := tagMap["default"]; ok && v != "" {
			defaultV = v
		}
		usage := ""
		if v, ok := tagMap["usage"]; ok && v != "" {
			usage = v
		}
		switch fv.Interface().(type) {
		case string:
			c.Flags().StringVarP(fv.Addr().Interface().(*string), name, short, defaultV, usage)
		case []string:
			c.Flags().StringArrayVarP(fv.Addr().Interface().(*[]string), name, short, strings.Split(defaultV, ","), usage)
		default:
			log.Fatal("unknown BuildCfg field type", log.Any("fieldName", f.Name), log.Any("fieldType", f.Type.String()))
		}
		if _, ok := tagMap["required"]; ok {
			err := c.MarkFlagRequired(name)
			if err != nil {
				log.Fatal("MarkFlagRequired error", log.Any("err", err))
			}
		}
	}
}
