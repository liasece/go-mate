package cmd

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/liasece/go-mate/config"
	ccontext "github.com/liasece/go-mate/context"
	"github.com/liasece/go-mate/gogen/writer"
	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func generateTmplToFile(ctx ccontext.ITmplContext, name string, toFile string, tmpl *config.TmplItem) {
	log.Debug(fmt.Sprintf("%s: generating %s", name, toFile))
	beginTime := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("%s: generated %s (%.2fs)", name, toFile, float64(time.Since(beginTime))/float64(time.Second)))
	}()

	if tmpl.OnlyCreate {
		notExists := false
		if _, err := os.Stat(toFile); errors.Is(err, os.ErrNotExist) {
			notExists = true
		} else if err != nil {
			log.Fatal("generateEntity tmpl check OnlyCreate os.Stat error", log.ErrorField(err))
			return
		}
		if !notExists {
			return
		}
	}
	c, err := ccontext.GetCodeFromTmpl(ctx, tmpl.From)
	if ctx.GetTerminate() {
		return
	}
	if err != nil {
		log.Fatal("generateEntity Tmpl GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err), log.Any("tmpl.From", tmpl.From))
		return
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
				log.Fatal("generateEntity Tmpl MergeGoFromFile error", log.ErrorField(err))
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
