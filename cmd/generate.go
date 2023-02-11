package cmd

import (
	"encoding/json"
	"regexp"

	"github.com/liasece/go-mate/config"
	ccontext "github.com/liasece/go-mate/context"
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
	coder_ast "github.com/liasece/gocoder/ast"
	"github.com/liasece/log"
)

type GenerateCfg struct {
	ConfigFile string `arg:"name: config file; short: f; usage: the config file path of target entity; required;"`
}

func newAstCoder(cfg *config.Config) (*coder_ast.CodeDecoder, error) {
	codePaths := append([]string{}, cfg.ImportGoCodePath...)
	for _, entityCfg := range cfg.Entity {
		entityCfg.CodeName = entityCfg.Name
		if entityCfg.EntityRealName != "" {
			entityRealName, err := utils.TemplateRaw(entityCfg.EntityRealName, &ccontext.ConfigTmplContext{
				VEntityName:  entityCfg.Name,
				VServiceName: entityCfg.Service,
			})
			if err != nil {
				log.Fatal("Generate EntityRealName TemplateRaw error", log.ErrorField(err))
				return nil, err
			}
			entityCfg.CodeName = entityRealName
		}
		entityPath, err := utils.TemplateRaw(entityCfg.EntityPath, &ccontext.ConfigTmplContext{
			VEntityName:  entityCfg.Name,
			VServiceName: entityCfg.Service,
		})
		if err != nil {
			log.Fatal("Generate TemplateRaw error", log.ErrorField(err))
			return nil, err
		}
		entityCfg.DecodedEntityPath = entityPath
		if entityCfg.Pkg == "" {
			entityCfg.Pkg = coder_ast.GetDirGoPackage(entityPath)
		}
		if entityCfg.Pkg == "" {
			if entityCfg.GetEnv("go", "mod") != "" {
				entityCfg.Pkg = entityCfg.GetEnv("go", "mod") + "/" + entityPath
			} else {
				entityCfg.Pkg = calGoFilePkgName(entityPath)
			}
		}
		codePaths = append(codePaths, entityPath)
	}
	astCoder, err := coder_ast.NewCodeDecoder(codePaths...)
	if err != nil {
		return nil, err
	}
	return astCoder, nil
}

func Generate(genCfg *GenerateCfg) {
	cfg, err := config.LoadConfig(genCfg.ConfigFile)
	if err != nil {
		log.Fatal("generate LoadConfig error", log.ErrorField(err), log.Any("genCfg", genCfg))
		return
	}
	err = log.InitLogByLevel(cfg.LogLevel)
	if err != nil {
		log.Fatal("generate InitLogByLevel error", log.ErrorField(err), log.Any("genCfg", genCfg))
		return
	}
	log.Debug("generate LoadConfig finish", log.Any("genCfg", genCfg), log.Any("cfg", cfg))
	{
		j, err := json.MarshalIndent(cfg.Entity, "", "\t")
		if err != nil {
			log.Fatal("Generate MarshalIndent error", log.ErrorField(err))
			return
		}
		log.Debug("entity generate config:\n" + string(j))
	}

	{
		// search entity names
		var appendEntity []*config.Entity
		astCoder, err := newAstCoder(cfg)
		if err != nil {
			log.Fatal("Generate NewCodeDecoder error", log.ErrorField(err))
			return
		}
		for i, entityCfg := range cfg.Entity {
			if entityCfg.Name == regexp.QuoteMeta(entityCfg.Name) {
				continue
			}
			switch entityCfg.EntityKind {
			case "methods":
				// log.Fatal("type name search not support in methods", log.ErrorField(err), log.Any("entityCfg", entityCfg))
				continue
			case "interface":
				// log.Fatal("type name search not support in interface", log.ErrorField(err), log.Any("entityCfg", entityCfg))
				continue
			default:
				entityNames := astCoder.SearchTypeNames(entityCfg.Pkg, entityCfg.Name)
				if len(entityNames) == 1 && entityNames[0] == entityCfg.Name {
					// no change
					continue
				} else {
					// replace this entity to newEntity
					log.Info("Generate appendEntity", log.Any("entityCfg.Pkg", entityCfg.Pkg), log.Any("entityNames", entityNames))
					cfg.Entity[i] = nil
					for _, entityName := range entityNames {
						newEntity := *entityCfg
						newEntity.Name = entityName
						appendEntity = append(appendEntity, &newEntity)
					}
				}
			}
		}
		newEntityList := make([]*config.Entity, 0, len(cfg.Entity)+len(appendEntity))
		for _, v := range cfg.Entity {
			if v != nil {
				newEntityList = append(newEntityList, v)
			}
		}
		newEntityList = append(newEntityList, appendEntity...)
		cfg.Entity = newEntityList
	}
	{
		// build entity info
		astCoder, err := newAstCoder(cfg)
		if err != nil {
			log.Fatal("Generate NewCodeDecoder error", log.ErrorField(err))
			return
		}
		for _, entityCfg := range cfg.Entity {
			entityCodeName := entityCfg.CodeName
			log.Debug("Generate begin", log.Any("entityCodeName", entityCodeName), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityPkg", entityCfg.Pkg), log.Any("entityKind", entityCfg.EntityKind))
			switch entityCfg.EntityKind {
			case "methods":
				// get entity type
				methods := astCoder.GetMethods(entityCfg.Pkg + "." + entityCodeName)
				if methods == nil {
					log.Debug("Generate methods LoadTypeFromSource not found", log.Any("entityCodeName", entityCodeName), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				cs := make([]gocoder.Codable, 0, len(methods))
				for _, v := range methods {
					cs = append(cs, v)
				}
				entityCfg.CodeType = gocoder.NewCode().C(cs...)
			case "interface":
				// get entity type
				interfaceType := astCoder.GetInterface(entityCfg.Pkg + "." + entityCodeName)
				if interfaceType == nil {
					log.Debug("Generate interface LoadTypeFromSource not found", log.Any("entityCodeName", entityCodeName), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				entityCfg.CodeType = interfaceType
			default:
				// get entity type
				entityType := astCoder.GetType(entityCfg.Pkg + "." + entityCodeName)
				if entityType == nil {
					log.Debug("Generate entity LoadTypeFromSource not found", log.Any("entityCodeName", entityCodeName), log.Any("entityCfg", entityCfg), log.Any("entityKind", entityCfg.EntityKind))
					continue
				}
				entityType.SetNamed(entityCodeName)
				{
					// debug info
					filedNames := make([]string, 0)
					for i := 0; i < entityType.NumField(); i++ {
						typ := entityType.Field(i).GetType()
						filedNames = append(filedNames, entityType.Field(i).GetName()+"("+typ.ShowString()+")")
					}
					log.Debug("Generate filedNames", log.Any("entityCodeName", entityCodeName), log.Any("entityFile", entityCfg.DecodedEntityPath), log.Any("entityPkg", entityCfg.Pkg),
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
