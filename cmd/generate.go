package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/liasece/go-mate/src/config"
	ccontext "github.com/liasece/go-mate/src/context"
	"github.com/liasece/go-mate/src/gogen/writer"
	"github.com/liasece/go-mate/src/gogen/writer/repo"
	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
)

type GenerateCfg struct {
	ConfigFile string `arg:"name: config file; short: f; usage: the config file path of target entity; required;"`
}

func Generate(genCfg *GenerateCfg) {
	cfg, err := config.LoadConfig(genCfg.ConfigFile)
	if err != nil {
		log.L(nil).Fatal("generate LoadConfig error", log.ErrorField(err), log.Any("genCfg", genCfg))
	}
	log.InitLogByLevel(cfg.LogLevel)
	log.L(nil).Debug("generate LoadConfig finish", log.Any("genCfg", genCfg), log.Any("cfg", cfg))
	{
		j, _ := json.MarshalIndent(cfg.Entity, "", "\t")
		log.L(nil).Debug("entity generate config:\n" + string(j))
	}

	{
		// build entity info
		codePaths := []string{}
		for _, entityCfg := range cfg.Entity {
			entityPath, err := utils.TemplateRaw(entityCfg.EntityPath, &ccontext.ConfigTmplContext{
				VEntityName:  entityCfg.Name,
				VServiceName: entityCfg.Service,
			})
			if err != nil {
				log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
				return
			}
			entityCfg.DecodedEntityPath = entityPath
			if entityCfg.Pkg == "" {
				if entityCfg.GetEnv("go", "mod") != "" {
					entityCfg.Pkg = entityCfg.GetEnv("go", "mod") + "/" + entityPath
				} else {
					entityCfg.Pkg = calGoFilePkgName(entityPath)
				}
			}
			codePaths = append(codePaths, entityPath)
		}
		astCoder, err := cde.NewASTCoder(codePaths...)
		if err != nil {
			log.Fatal("generateEntity NewASTCoder error", log.ErrorField(err))
			return
		}
		for _, entityCfg := range cfg.Entity {
			// get entity type
			entityType, err := astCoder.GetType(entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
			if err != nil {
				log.Fatal("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityCfg.Name", entityCfg.Name))
				return
			}
			if entityType == nil {
				log.Fatal("generateEntity LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg))
				return
			}
			entityType.SetNamed(entityCfg.Name)
			{
				filedNames := make([]string, 0)
				for i := 0; i < entityType.NumField(); i++ {
					typ := entityType.Field(i).GetType()
					filedNames = append(filedNames, entityType.Field(i).GetName()+"("+typ.ShowString()+")")
				}
				log.Debug("generateEntity filedNames", log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("entityNames", entityCfg.Name),
					log.Any("entityCfg.Name", entityCfg.Name), log.Any("filedNames", filedNames))
			}
			entityCfg.CodeType = entityType
		}
	}

	for _, entity := range cfg.Entity {
		generateEntity(entity)
	}
}

func generateEntity(entityCfg *config.Entity) {
	enGameEntry := repo.NewRepositoryWriterByType(entityCfg.CodeType, entityCfg.Name, entityCfg.Pkg, entityCfg.Service, "", "", "", "", "")
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

	if entityCfg.OptFilePath != "" {
		optCode := gocoder.NewCode()
		optCode.C(enGameEntry.GetFilterTypeCode(), enGameEntry.GetUpdaterTypeCode(), enGameEntry.GetSorterTypeCode())
		if entityCfg.NoSelector == nil || !*entityCfg.NoSelector {
			optCode.C(enGameEntry.GetSelectorTypeCode())
		}

		optFile, err := utils.TemplateRaw(entityCfg.OptFilePath, ccontext.NewEntityTmplContext(enGameEntry))
		if err != nil {
			log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
			return
		}
		if optFile != "" {
			log.Info(fmt.Sprintf("%s: generating %s", entityCfg.Name, entityCfg.OptFilePath))
			optFileDirPath := filepath.Dir(optFile)
			optPkg := entityCfg.GetEnv("go", "mod") + "/" + optFileDirPath
			if optPkg == "" {
				optPkg = entityCfg.Pkg
			}
			optPkgName := filepath.Base(optPkg)
			// log.L(nil).Fatal("entityCfg.OptFilePath gen", log.Any("optPkg", optPkg))
			err := gocoder.WriteToFile(optFile, optCode, gocoder.NewToCodeOpt().PkgName(optPkgName).PkgPath(optPkg))
			if err != nil {
				log.L(nil).Fatal("generateEntity optFile WriteToFile error", log.ErrorField(err), log.Any("optFile", optFile))
			}
		}
	}
}

func generateEntityProtoType(entityCfg *config.Entity, enGameEntry *repo.RepositoryWriter) {
	protoTypeFile, err := utils.TemplateRaw(entityCfg.ProtoTypeFile, ccontext.NewEntityTmplContext(enGameEntry))
	if err != nil {
		log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
		return
	}
	if protoTypeFile != "" {
		log.Debug(fmt.Sprintf("%s: generating %s", entityCfg.Name, protoTypeFile))
		beginTime := time.Now()
		defer func() {
			log.Info(fmt.Sprintf("%s: generated %s (%.2fs)", entityCfg.Name, protoTypeFile, float64(time.Now().Sub(beginTime))/float64(time.Second)))
		}()
		ts := []gocoder.Type{
			entityCfg.CodeType,
			enGameEntry.Filter.GetType(),
			enGameEntry.Updater.GetType(),
			enGameEntry.Sorter.GetType(),
		}
		if entityCfg.NoSelector == nil || !*entityCfg.NoSelector {
			ts = append(ts, enGameEntry.Selector.GetType())
		}
		writer.StructToProto(protoTypeFile, entityCfg.ProtoTypeFileIndent, ts...)
	}
}

func generateEntityTmplItem(entityCfg *config.Entity, enGameEntry *repo.RepositoryWriter, tmpl *config.TmplItem) {
	toFile, err := utils.TemplateRaw(tmpl.To, ccontext.NewTmplContext(enGameEntry, tmpl))
	if err != nil {
		log.Fatal("generateEntity TemplateRaw error", log.ErrorField(err))
		return
	}
	log.Debug(fmt.Sprintf("%s: generating %s", entityCfg.Name, toFile))
	beginTime := time.Now()
	defer func() {
		log.Info(fmt.Sprintf("%s: generated %s (%.2fs)", entityCfg.Name, toFile, float64(time.Now().Sub(beginTime))/float64(time.Second)))
	}()

	if tmpl.OnlyCreate {
		notExists := false
		if _, err := os.Stat(toFile); errors.Is(err, os.ErrNotExist) {
			notExists = true
		} else if err != nil {
			log.L(nil).Fatal("generateEntity tmpl check OnlyCreate os.Stat error", log.ErrorField(err))
			return
		}
		if !notExists {
			return
		}
	}
	tmplCtx := ccontext.NewTmplContext(enGameEntry, tmpl)
	c, err := ccontext.GetEntityRepositoryCodeFromTmpl(enGameEntry, tmpl.From, tmplCtx)
	if tmplCtx.GetTerminate() {
		return
	}
	if err != nil {
		log.Fatal("generateEntity Tmpl GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err), log.Any("tmpl.From", tmpl.From))
		return
	} else {
		if tmpl.Merge {
			switch tmpl.Type {
			case config.TmplItemTypeProto:
				writer.MergeProtoFromFile(toFile, gocoder.ToCode(c, gocoder.NewToCodeOpt().PkgName("")))
			case config.TmplItemTypeGo:
				writer.MergeGoFromFile(toFile, gocoder.ToCode(c, gocoder.NewToCodeOpt().PkgName("")))
			case config.TmplItemTypeGraphQL:
				writer.MergeGraphQLFromFile(toFile, gocoder.ToCode(c, gocoder.NewToCodeOpt().PkgName("")))
			default:
				log.Fatal("generateEntity Template merge type not support", log.Any("tmpl", tmpl))
			}
		} else {
			err := gocoder.WriteToFile(toFile, c, gocoder.NewToCodeOpt().PkgName(""))
			if err != nil {
				log.L(nil).Fatal("generateEntity tmpl WriteToFile error", log.ErrorField(err), log.Any("toFile", toFile))
			}
		}
	}
}
