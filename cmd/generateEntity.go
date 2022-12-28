package cmd

import (
	"fmt"
	"time"

	"github.com/liasece/go-mate/src/config"
	ccontext "github.com/liasece/go-mate/src/context"
	"github.com/liasece/go-mate/src/gogen/writer"
	"github.com/liasece/go-mate/src/gogen/writer/repo"
	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func generateEntity(entityCfg *config.Entity) {
	if entityCfg.CodeType == nil {
		return
	}
	entityType := entityCfg.CodeType.(gocoder.Type)
	enGameEntry := repo.NewRepositoryWriterByType(entityType, entityCfg.Name, entityCfg.Pkg, entityCfg.Service, "", "", "", "", "")
	enGameEntry.EntityCfg = entityCfg
	enGameEntry.Filter, _ = enGameEntry.GetFilterTypeStructCode()
	enGameEntry.Updater, _ = enGameEntry.GetUpdaterTypeStructCode()
	enGameEntry.Sorter, _ = enGameEntry.GetSorterTypeStructCode()
	enGameEntry.Selector, _ = enGameEntry.GetSelectorTypeStructCode()

	{
		generateEntityProtoType(entityCfg, enGameEntry)
	}

	for _, tmpl := range entityCfg.Tmpl {
		generateEntityTmplItem(entityCfg, enGameEntry, tmpl)
	}
}

func generateEntityProtoType(entityCfg *config.Entity, enGameEntry *repo.RepositoryWriter) {
	protoTypeFile, err := utils.TemplateRaw(entityCfg.ProtoTypeFile, ccontext.NewEntityTmplContext(nil, enGameEntry))
	if err != nil {
		log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
		return
	}
	if protoTypeFile != "" {
		log.Debug(fmt.Sprintf("%s: generating %s", entityCfg.Name, protoTypeFile))
		beginTime := time.Now()
		defer func() {
			log.Info(fmt.Sprintf("%s: generated %s (%.2fs)", entityCfg.Name, protoTypeFile, float64(time.Since(beginTime))/float64(time.Second)))
		}()
		var updater gocoder.Struct
		{
			// build pb updater
			updater, _ = enGameEntry.GetUpdaterTypeStructCode(repo.TypeOptAddSlicePBEmpty(true))
		}
		ts := []gocoder.Type{
			entityCfg.CodeType.(gocoder.Type),
			enGameEntry.Filter.GetType(),
			updater.GetType(),
			enGameEntry.Sorter.GetType(),
		}
		if entityCfg.NoSelector == nil || !*entityCfg.NoSelector {
			ts = append(ts, enGameEntry.Selector.GetType())
		}
		err = writer.StructToProto(protoTypeFile, entityCfg.ProtoTypeFileIndent, ts...)
		if err != nil {
			log.Fatal("generateEntity StructToProto error", log.ErrorField(err))
			return
		}
	}
}

func generateEntityTmplItem(entityCfg *config.Entity, enGameEntry *repo.RepositoryWriter, tmpl *config.TmplItem) {
	toFile, err := utils.TemplateRaw(tmpl.To, ccontext.NewEntityTmplContext(ccontext.NewTmplContext(tmpl, entityCfg), enGameEntry))
	if err != nil {
		log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
		return
	}
	tmplCtx := ccontext.NewEntityTmplContext(ccontext.NewTmplContext(tmpl, entityCfg), enGameEntry)
	generateTmplToFile(tmplCtx, entityCfg.Name, toFile, tmpl)
}
