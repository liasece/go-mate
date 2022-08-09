package main

import (
	"encoding/json"
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
	path, _ := filepath.Split(entityCfg.EntityPath)
	if entityCfg.Pkg == "" {
		entityCfg.Pkg = calGoFilePkgName(entityCfg.EntityPath)
	}
	log.Info("buildRunner begin", log.Any("entityFile", entityCfg.EntityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("path", path), log.Any("entityNames", entityCfg.Name))
	// log.Info("buildRunner begin", log.Any("entityFile", entityCfg.EntityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("path", path), log.Any("entityNames", entityCfg.Name), log.Any("entityCfg", entityCfg))

	optCode := gocoder.NewCode()
	repositoryInterfaceCode := gocoder.NewCode()

	t, err := cde.LoadTypeFromSource(entityCfg.EntityPath, entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
	if err != nil {
		log.Error("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityCfg.EntityPath), log.Any("entityCfg.Name", entityCfg.Name))
		return
	}
	if t == nil {
		log.Error("buildRunner LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg))
		return
	}
	t.SetNamed(entityCfg.Name)
	{
		filedNames := make([]string, 0)
		for i := 0; i < t.NumField(); i++ {
			typ := t.Field(i).GetType()
			filedNames = append(filedNames, t.Field(i).GetName()+"("+typ.ShowString()+")")
		}
		log.Info("buildRunner filedNames", log.Any("entityCfg.Name", entityCfg.Name), log.Any("filedNames", filedNames))
	}
	enGameEntry := repo.NewRepositoryWriterByType(t, entityCfg.Name, entityCfg.Pkg, entityCfg.Service, "", "", "", "")
	enGameEntry.EntityCfg = entityCfg
	optCode.C(enGameEntry.GetFilterTypeCode(), enGameEntry.GetUpdaterTypeCode(), enGameEntry.GetSorterTypeCode())
	repositoryInterfaceCode.C(enGameEntry.GetEntityRepositoryInterfaceCode())

	if entityCfg.ProtoTypeFile != "" {
		writer.StructToProto(entityCfg.ProtoTypeFile, t, entityCfg.ProtoTypeFileIndent)
		filterStr, _ := enGameEntry.GetFilterTypeStructCode()
		writer.StructToProto(entityCfg.ProtoTypeFile, filterStr.GetType(), entityCfg.ProtoTypeFileIndent)
		updaterStr, _ := enGameEntry.GetUpdaterTypeStructCode()
		writer.StructToProto(entityCfg.ProtoTypeFile, updaterStr.GetType(), entityCfg.ProtoTypeFileIndent)
		sorterStr, _ := enGameEntry.GetSorterTypeStructCode()
		writer.StructToProto(entityCfg.ProtoTypeFile, sorterStr.GetType(), entityCfg.ProtoTypeFileIndent)
	}

	for _, tmpl := range entityCfg.Tmpl {
		toFile, err := gocoder.TemplateRaw(tmpl.To, enGameEntry.NewTmplRepositoryEnv(), nil)
		if err != nil {
			log.Error("buildRunner TemplateRaw error", log.ErrorField(err))
			return
		}
		c, err := enGameEntry.GetEntityRepositoryCodeFromTmpl(tmpl.From)
		if err != nil {
			log.Error("buildRunner Tmpl GetEntityRepositoryCodeFromTmpl error", log.ErrorField(err), log.Any("tmpl.From", tmpl.From))
		} else {
			if tmpl.Merge {
				writer.MergeProtoFromFile(toFile, gocoder.ToCode(c, gocoder.NewToCodeOpt().PkgName("")))
			} else {
				gocoder.WriteToFile(toFile, c, gocoder.NewToCodeOpt().PkgName(""))
			}
		}
	}

	if entityCfg.CopierFile != "" {
		var info *writer.ProtoInfo
		if entityCfg.ProtoTypeFile != "" {
			info, _ = writer.ReadProtoInfo(entityCfg.ProtoTypeFile)
		}
		if info != nil {
			optPkg := pkgInReference(entityCfg.Pkg)
			if entityCfg.EntityOptPkg != "" {
				optPkg = pkgInReference(entityCfg.EntityOptPkg)
			}
			entityPkg := pkgInReference(entityCfg.Pkg)
			infoPkg := pkgInReference(info.Package)
			var names [][2]string = [][2]string{
				{entityPkg + "." + entityCfg.Name, infoPkg + entityCfg.OutputCopierProtoPkgSuffix + "." + entityCfg.Name},
				{infoPkg + entityCfg.OutputCopierProtoPkgSuffix + "." + entityCfg.Name, entityPkg + "." + entityCfg.Name},
				{infoPkg + entityCfg.OutputCopierProtoPkgSuffix + "." + enGameEntry.GetFilterTypeStructCodeStruct().GetName(), optPkg + "." + enGameEntry.GetFilterTypeStructCodeStruct().GetName()},
				{infoPkg + entityCfg.OutputCopierProtoPkgSuffix + "." + enGameEntry.GetUpdaterTypeStructCodeStruct().GetName(), optPkg + "." + enGameEntry.GetUpdaterTypeStructCodeStruct().GetName()},
				{infoPkg + entityCfg.OutputCopierProtoPkgSuffix + "." + enGameEntry.GetSorterTypeStructCodeStruct().GetName(), optPkg + "." + enGameEntry.GetSorterTypeStructCodeStruct().GetName()},
			}
			writer.FillCopierLine(entityCfg.CopierFile, names)
		}
	}

	if entityCfg.EntityPath != "" {
		if entityCfg.Pkg == "" {
			entityCfg.Pkg = calGoFilePkgName(entityCfg.EntityPath)
		}
		gocoder.WriteToFile(entityCfg.EntityPath, optCode, gocoder.NewToCodeOpt().PkgName(entityCfg.Pkg))
	}
}
