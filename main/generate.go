package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/liasece/go-mate/src/config"
	"github.com/liasece/go-mate/src/gogen/writer"
	"github.com/liasece/go-mate/src/gogen/writer/repo"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/log"
)

type GenerateCfg struct {
	ConfigFile string `arg:"name: config file; short: f; usage: the config file path of target entity; required;"`
}

func generate(genCfg *GenerateCfg) {
	cfg, err := config.LoadConfig(genCfg.ConfigFile)
	if err != nil {
		log.L(nil).Fatal("generate LoadConfig error", log.ErrorField(err), log.Any("genCfg", genCfg))
	}
	log.L(nil).Debug("generate LoadConfig finish", log.Any("genCfg", genCfg), log.Any("cfg", cfg))
	j, _ := json.MarshalIndent(cfg, "", "\t")
	log.L(nil).Debug("generate config:\n" + string(j))

	for _, entity := range cfg.Entity {
		generateEntity(entity)
	}
}

func generateEntity(entityCfg *config.Entity) {
	tmpPaths := make([]string, 0)
	_ = tmpPaths
	tmpFiles := make([]string, 0)
	_ = tmpFiles
	entityPath, err := gocoder.TemplateRaw(entityCfg.EntityPath, &config.ConfigTmplContext{
		VEntityName:  entityCfg.Name,
		VServiceName: entityCfg.Service,
	}, nil)
	if err != nil {
		log.Error("generateEntity TemplateRaw error", log.ErrorField(err))
		return
	}
	path, _ := filepath.Split(entityPath)
	if entityCfg.Pkg == "" {
		if entityCfg.GetEnv("go", "mod") != "" {
			entityCfg.Pkg = entityCfg.GetEnv("go", "mod") + "/" + entityPath
		} else {
			entityCfg.Pkg = calGoFilePkgName(entityPath)
		}
	}
	log.Debug("generateEntity begin", log.Any("entityFile", entityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("path", path), log.Any("entityNames", entityCfg.Name))
	// log.Info("generateEntity begin", log.Any("entityFile", entityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("path", path), log.Any("entityNames", entityCfg.Name), log.Any("entityCfg", entityCfg))

	entityType, err := cde.LoadTypeFromSource(entityPath, entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
	if err != nil {
		log.Error("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityPath), log.Any("entityCfg.Name", entityCfg.Name))
		return
	}
	if entityType == nil {
		log.Error("generateEntity LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg))
		return
	}
	entityType.SetNamed(entityCfg.Name)
	{
		filedNames := make([]string, 0)
		for i := 0; i < entityType.NumField(); i++ {
			typ := entityType.Field(i).GetType()
			filedNames = append(filedNames, entityType.Field(i).GetName()+"("+typ.ShowString()+")")
		}
		log.Info("generateEntity filedNames", log.Any("entityFile", entityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("path", path), log.Any("entityNames", entityCfg.Name),
			log.Any("entityCfg.Name", entityCfg.Name), log.Any("filedNames", filedNames))
	}
	enGameEntry := repo.NewRepositoryWriterByType(entityType, entityCfg.Name, entityCfg.Pkg, entityCfg.Service, "", "", "", "")
	enGameEntry.EntityCfg = entityCfg

	{
		protoTypeFile, err := gocoder.TemplateRaw(entityCfg.ProtoTypeFile, enGameEntry.NewEntityTmplContext(), nil)
		if err != nil {
			log.Error("generateEntity TemplateRaw error", log.ErrorField(err))
			return
		}
		if protoTypeFile != "" {
			writer.StructToProto(protoTypeFile, entityType, entityCfg.ProtoTypeFileIndent)
			filterStr, _ := enGameEntry.GetFilterTypeStructCode()
			enGameEntry.Filter = filterStr
			writer.StructToProto(protoTypeFile, filterStr.GetType(), entityCfg.ProtoTypeFileIndent)
			updaterStr, _ := enGameEntry.GetUpdaterTypeStructCode()
			enGameEntry.Updater = updaterStr
			writer.StructToProto(protoTypeFile, updaterStr.GetType(), entityCfg.ProtoTypeFileIndent)
			sorterStr, _ := enGameEntry.GetSorterTypeStructCode()
			enGameEntry.Sorter = sorterStr
			writer.StructToProto(protoTypeFile, sorterStr.GetType(), entityCfg.ProtoTypeFileIndent)
		}
	}

	for _, tmpl := range entityCfg.Tmpl {
		toFile, err := gocoder.TemplateRaw(tmpl.To, enGameEntry.NewEntityTmplContext(), nil)
		if err != nil {
			log.Error("generateEntity TemplateRaw error", log.ErrorField(err))
			return
		}
		if tmpl.OnlyCreate {
			notExists := false
			if _, err := os.Stat(toFile); errors.Is(err, os.ErrNotExist) {
				notExists = true
			} else if err != nil {
				log.L(nil).Fatal("generateEntity tmpl check OnlyCreate os.Stat error", log.ErrorField(err))
				continue
			}
			if !notExists {
				continue
			}
		}
		c, err := enGameEntry.GetEntityRepositoryCodeFromTmpl(tmpl.From)
		if err != nil {
			log.Error("generateEntity Tmpl GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err), log.Any("tmpl.From", tmpl.From))
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
					log.Error("generateEntity Template merge type not support", log.Any("tmpl", tmpl))
				}
			} else {
				gocoder.WriteToFile(toFile, c, gocoder.NewToCodeOpt().PkgName(""))
			}
		}
	}

	if entityCfg.OptFilePath != "" {
		optCode := gocoder.NewCode()
		optCode.C(enGameEntry.GetFilterTypeCode(), enGameEntry.GetUpdaterTypeCode(), enGameEntry.GetSorterTypeCode())

		optFile, err := gocoder.TemplateRaw(entityCfg.OptFilePath, enGameEntry.NewEntityTmplContext(), nil)
		if err != nil {
			log.Error("generateEntity TemplateRaw error", log.ErrorField(err))
			return
		}
		if optFile != "" {
			optPkg := calGoFilePkgName(optFile)
			if optPkg == "" {
				optPkg = entityCfg.Pkg
			}
			gocoder.WriteToFile(optFile, optCode, gocoder.NewToCodeOpt().PkgName(optPkg))
		}
	}
}
