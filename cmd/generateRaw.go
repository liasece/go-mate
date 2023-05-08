package cmd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/liasece/go-mate/config"
	ccontext "github.com/liasece/go-mate/context"
	"github.com/liasece/log"
)

func generateRaw(raw *config.Raw) {
	log.Debug(fmt.Sprintf("raw: generating %+v", raw))
	beginTime := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("raw: generated %+v (%.2fs)", raw, float64(time.Since(beginTime))/float64(time.Second)))
	}()

	configs := make([][]map[string]interface{}, 0)
	{
		entries, err := filepath.Glob(raw.Config)
		if err != nil {
			log.Fatal("generateRaw read config dir error", log.ErrorField(err), log.String("dir", raw.Config))
		}
		log.Info("generateRaw read config dir", log.Any("config", raw.Config), log.Any("entries", entries))
		for _, entry := range entries {
			var value []map[string]interface{}
			err := config.LoadConfigTo(entry, &value)
			if err != nil {
				log.Fatal("generateRaw load config error", log.ErrorField(err), log.String("entry", entry))
			}
			configs = append(configs, value)
		}
	}

	for _, tmpl := range raw.Tmpl {
		for _, c := range configs {
			for _, v := range c {
				tmplCtx := ccontext.NewRawTmplContext(ccontext.NewTmplContext(tmpl, nil), v)
				generateTmplToFile(tmplCtx, tmpl)
			}
		}
	}
}
