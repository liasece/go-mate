package cmd

import (
	"encoding/json"

	"github.com/liasece/go-mate/src/config"
	ccontext "github.com/liasece/go-mate/src/context"
	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
	coder_ast "github.com/liasece/gocoder/ast"
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
		astCoder, err := coder_ast.NewASTCoder(codePaths...)
		if err != nil {
			log.Fatal("generateEntity NewASTCoder error", log.ErrorField(err))
			return
		}
		for _, entityCfg := range cfg.Entity {
			switch entityCfg.EntityKind {
			case "methods":
				// get entity type
				methods, err := astCoder.GetMethods(entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
				if err != nil {
					log.Fatal("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityKind", entityCfg.EntityKind))
					return
				}
				if methods == nil {
					log.Debug("generateEntity methods LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				cs := make([]gocoder.Codable, 0, len(methods))
				for _, v := range methods {
					cs = append(cs, v)
				}
				entityCfg.CodeType = gocoder.NewCode().C(cs...)
			case "interface":
				// get entity type
				interfaceType, err := astCoder.GetInterface(entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
				if err != nil {
					log.Fatal("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityKind", entityCfg.EntityKind))
					return
				}
				if interfaceType == nil {
					log.Debug("generateEntity interface LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				entityCfg.CodeType = interfaceType
			default:
				// get entity type
				entityType, err := astCoder.GetType(entityCfg.Name, gocoder.NewToCodeOpt().PkgPath(entityCfg.Pkg))
				if err != nil {
					log.Fatal("LoadTypeFromSource error", log.ErrorField(err), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityKind", entityCfg.EntityKind))
					return
				}
				if entityType == nil {
					log.Debug("generateEntity entity LoadTypeFromSource not found", log.Any("entityCfg.Name", entityCfg.Name), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				entityType.SetNamed(entityCfg.Name)
				{
					// debug info
					filedNames := make([]string, 0)
					for i := 0; i < entityType.NumField(); i++ {
						typ := entityType.Field(i).GetType()
						filedNames = append(filedNames, entityType.Field(i).GetName()+"("+typ.ShowString()+")")
					}
					log.Debug("generateEntity filedNames", log.Any("entityName", entityCfg.Name), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityPkg", entityCfg.Pkg),
						log.Any("filedNames", filedNames), log.Any("entityKind", entityCfg.EntityKind))
				}
				entityCfg.CodeType = entityType
			}
		}
	}

	for _, entityCfg := range cfg.Entity {
		switch entityCfg.EntityKind {
		case "methods":
			generateMethods(entityCfg)
		default:
			generateEntity(entityCfg)
		}
	}
}
