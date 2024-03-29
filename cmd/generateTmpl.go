package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/liasece/go-mate/config"
	ccontext "github.com/liasece/go-mate/context"
	"github.com/liasece/go-mate/gogen/writer"
	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func generateTmplCheck(tmpl *config.TmplItem, toFile string) bool {
	if tmpl.OnlyCreate {
		fileExists := false
		toFileName := filepath.Base(toFile)
		files, err := os.ReadDir(filepath.Dir(toFile))
		if err != nil {
			log.Fatal("generateEntity tmpl check OnlyCreate os.Stat error", log.ErrorField(err))
			return false
		}
		for _, f := range files {
			if strings.EqualFold(f.Name(), toFileName) {
				fileExists = true
				break
			}
		}
		if fileExists {
			return false
		}
	}
	return true
}

func generateEntityTmplToFile(ctx ccontext.ITmplContext, name string, fromFile string, toFile string, tmpl *config.TmplItem) {
	log.Debug(fmt.Sprintf("%s: generating %s", name, toFile))
	beginTime := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("%s: generated %s (%.2fs)", name, toFile, float64(time.Since(beginTime))/float64(time.Second)))
	}()
	ctx.FromFilePath(fromFile)
	ctx.ToFilePath(toFile)
	generateTmplToFile(ctx, tmpl)
}

func generateTmplToFile(ctx ccontext.ITmplContext, tmpl *config.TmplItem) {
	toFile := ctx.GetToFilePath()
	if toFile != "" && !generateTmplCheck(tmpl, toFile) {
		return
	}
	fromFile := ctx.GetFromFilePath()
	if fromFile == "" {
		log.Fatal("generateEntity Tmpl GetEntityRepositoryCodeFromTmpl error: fromFile not set", log.Any("fromFile", fromFile))
		return
	}
	c, err := ccontext.GetCodeFromTmpl(ctx, fromFile)
	if ctx.GetTerminate() {
		return
	}
	if err != nil {
		log.Fatal("generateEntity Tmpl GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err), log.Any("fromFile", fromFile))
		return
	}
	if ctx.GetToFilePath() != "" && ctx.GetToFilePath() != toFile && !generateTmplCheck(tmpl, toFile) {
		return
	}
	toFile = ctx.GetToFilePath()
	if toFile == "" {
		log.Fatal("generateEntity Tmpl GetToFilePath toFile == \"\"", log.Any("tmpl", tmpl))
		return
	}
	if tmpl.Type == "" {
		switch {
		case strings.HasSuffix(toFile, ".proto"):
			tmpl.Type = config.TmplItemTypeProto
		case strings.HasSuffix(toFile, ".go"):
			tmpl.Type = config.TmplItemTypeGo
		case strings.HasSuffix(toFile, ".graphql"):
			tmpl.Type = config.TmplItemTypeGraphQL
		default:
			log.Fatal("generateEntity Tmpl Type == \"\"", log.Any("tmpl", tmpl))
			return
		}
	}
	if tmpl.Merge {
		codeStr := gocoder.ToCode(c, gocoder.NewToCodeOpt().PkgName(""))
		switch tmpl.Type {
		case config.TmplItemTypeProto:
			err := writer.MergeProtoFromFile(toFile, codeStr)
			if err != nil {
				log.Fatal("generateEntity Tmpl MergeProtoFromFile error", log.ErrorField(err))
				return
			}
		case config.TmplItemTypeGo:
			err := writer.MergeGoFromFile(toFile, codeStr)
			if err != nil {
				log.Fatal("generateEntity Tmpl MergeGoFromFile error", log.ErrorField(err), log.Any("codeStr", codeStr))
				return
			}
		case config.TmplItemTypeGraphQL:
			err := writer.MergeGraphQLFromFile(toFile, codeStr)
			if err != nil {
				log.Fatal("generateEntity Tmpl MergeGraphQLFromFile error", log.ErrorField(err))
				return
			}
		default:
			log.Fatal("generateEntity Template merge type not support", log.Any("tmpl", tmpl))
		}
	} else {
		err := gocoder.WriteToFile(toFile, c, gocoder.NewToCodeOpt().PkgName(""))
		if err != nil {
			log.Fatal("generateEntity tmpl WriteToFile error", log.ErrorField(err), log.Any("toFile", toFile))
		}
	}
}
