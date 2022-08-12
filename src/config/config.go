package config

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/liasece/log"
	"gopkg.in/yaml.v2"
)

type TmplItemType string

const (
	TmplItemTypeGo      TmplItemType = "go"
	TmplItemTypeProto   TmplItemType = "proto"
	TmplItemTypeGraphQL TmplItemType = "graphql"
)

type TmplItem struct {
	From       string       `json:"from" yaml:"from"`
	To         string       `json:"to" yaml:"to"`
	Type       TmplItemType `json:"type,omitempty" yaml:"type,omitempty"`
	Merge      bool         `json:"merge,omitempty" yaml:"merge,omitempty"`
	OnlyCreate bool         `json:"onlyCreate,omitempty" yaml:"onlyCreate,omitempty"`
}

type ServiceBase struct {
	EntityPath                 string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	ProtoTypeFile              string                       `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	ProtoTypeFileIndent        string                       `json:"protoTypeFileIndent,omitempty" yaml:"protoTypeFileIndent,omitempty"`
	EntityOptPkg               string                       `json:"entityOptPkg,omitempty" yaml:"entityOptPkg,omitempty"`
	OutputCopierProtoPkgSuffix string                       `json:"outputCopierProtoPkgSuffix,omitempty" yaml:"outputCopierProtoPkgSuffix,omitempty"`
	Env                        map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
}

type Service struct {
	Name        string `json:"name" yaml:"name"`
	ServiceBase `json:",inline" yaml:",inline"`
}

type Base struct {
	Service map[string]*Service `json:"service" yaml:"service"`
}

type Comment struct {
	Doc string `json:"doc,omitempty" yaml:"doc,omitempty"`
}

type Config struct {
	Comment               `json:",inline" yaml:",inline"`
	Base                  `json:"base" yaml:"base"`
	Entity                []*Entity                    `json:"entity,omitempty" yaml:"entity,omitempty"`
	EntityPrefab          map[string]*EntityPrefab     `json:"entityPrefab,omitempty" yaml:"entityPrefab,omitempty"`
	BuildEntityWithPrefab map[string][]string          `json:"buildEntityWithPrefab,omitempty" yaml:"buildEntityWithPrefab,omitempty"`
	Env                   map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
}

type EntityPrefab struct {
	Comment       `json:",inline" yaml:",inline"`
	Name          string                       `json:"name" yaml:"name"`
	Pkg           string                       `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields        []*EntityField               `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service       string                       `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg    string                       `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	Prefab        []string                     `json:"prefab,omitempty" yaml:"prefab,omitempty"`
	Tmpl          []*TmplItem                  `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
	Env           map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	EntityPath    string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	ProtoTypeFile string                       `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	OptFilePath   string                       `json:"optFilePath,omitempty" yaml:"optFilePath,omitempty"`
}

type Entity struct {
	Comment                    `json:",inline" yaml:",inline"`
	Name                       string                       `json:"name" yaml:"name"`
	Pkg                        string                       `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields                     []*EntityField               `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service                    string                       `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg                 string                       `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	Prefab                     []string                     `json:"prefab,omitempty" yaml:"prefab,omitempty"`
	Tmpl                       []*TmplItem                  `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
	Env                        map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	ProtoTypeFileIndent        string                       `json:"protoTypeFileIndent,omitempty" yaml:"protoTypeFileIndent,omitempty"`
	EntityPath                 string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	ProtoTypeFile              string                       `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	EntityOptPkg               string                       `json:"entityOptPkg,omitempty" yaml:"entityOptPkg,omitempty"`
	OutputCopierProtoPkgSuffix string                       `json:"outputCopierProtoPkgSuffix,omitempty" yaml:"outputCopierProtoPkgSuffix,omitempty"`
	OptFilePath                string                       `json:"optFilePath,omitempty" yaml:"optFilePath,omitempty"`
}

type EntityFieldType string

const (
	EntityFieldTypeStruct EntityFieldType = "struct"

	EntityFieldTypeString EntityFieldType = "string"
	EntityFieldTypeTime   EntityFieldType = "time"
	EntityFieldTypeBool   EntityFieldType = "bool"
	EntityFieldTypeInt    EntityFieldType = "int"
	EntityFieldTypeByte   EntityFieldType = "byte"
	EntityFieldTypeInt8   EntityFieldType = "int8"
	EntityFieldTypeInt16  EntityFieldType = "int16"
	EntityFieldTypeInt32  EntityFieldType = "int32"
	EntityFieldTypeInt64  EntityFieldType = "int64"
	EntityFieldTypeUint8  EntityFieldType = "uint8"
	EntityFieldTypeUint16 EntityFieldType = "uint16"
	EntityFieldTypeUint32 EntityFieldType = "uint32"
	EntityFieldTypeUint64 EntityFieldType = "uint64"
)

type EntityField struct {
	Comment `json:",inline" yaml:",inline"`
	Name    string          `json:"name" yaml:"name"`
	Type    EntityFieldType `json:"type" yaml:"type"`
	Entity  *Entity         `json:"entity" yaml:"entity"`
}

func (c *TmplItem) AfterLoad() {
	if c.Type == "" {
		if strings.HasSuffix(c.To, ".go") {
			c.Type = TmplItemTypeGo
		} else if strings.HasSuffix(c.To, ".proto") {
			c.Type = TmplItemTypeProto
		} else if strings.HasSuffix(c.To, ".graphql") {
			c.Type = TmplItemTypeGraphQL
		} else {
			log.L(nil).Fatal("TmplItem AfterLoad unknown tmpl type", log.Any("tmpl", c))
		}
	}
}

func (c *ServiceBase) AfterLoad() {
	getIndent := func() string {
		if c.ProtoTypeFileIndent == "" {
			return "\t"
		}
		if strings.HasPrefix(c.ProtoTypeFileIndent, "$") {
			switch c.ProtoTypeFileIndent[1:] {
			case "tab":
				return "\t"
			default:
				a, _ := strconv.Atoi(c.ProtoTypeFileIndent[1:])
				return strings.Repeat(" ", a)
			}
		}
		return c.ProtoTypeFileIndent
	}
	c.ProtoTypeFileIndent = getIndent()
}

func (c *Config) AfterLoad() {
	// build prefab
	for _, prefab := range c.EntityPrefab {
		prefab.Env = mergeEnv(prefab.Env, c.Env)
	}

	for _, prefab := range c.EntityPrefab {
		for _, innerPrefab := range prefab.Prefab {
			c.EntityPrefab[innerPrefab].ApplyToPrefab(prefab)
		}
		for _, tmpl := range prefab.Tmpl {
			tmpl.AfterLoad()
		}
	}

	for prefabName, entityNameList := range c.BuildEntityWithPrefab {
		for _, entityName := range entityNameList {
			find := false
			for _, v := range c.Entity {
				if v.Name == entityName {
					find = true
					break
				}
			}
			if !find {
				c.Entity = append(c.Entity, &Entity{
					Name:   entityName,
					Prefab: []string{prefabName},
				})
			} else {
				log.L(nil).Fatal("Config AfterLoad build prefab target entity already exists", log.Any("entityName", entityName))
			}
		}
	}

	for k, service := range c.Service {
		service.Name = k
		service.ServiceBase.AfterLoad()
		service.Env = mergeEnv(service.Env, c.Env)
	}
	for _, entity := range c.Entity {
		for _, prefab := range entity.Prefab {
			if prefab, ok := c.EntityPrefab[prefab]; ok {
				prefab.ApplyToEntity(entity)
			} else {
				log.L(nil).Fatal("Config AfterLoad entity prefab not found", log.Any("entity", entity))
			}
		}
		if entity.Service != "" {
			if service, ok := c.Service[entity.Service]; ok {
				service.ApplyToEntity(entity)
			} else {
				log.L(nil).Fatal("Config AfterLoad service not found", log.Any("entity", entity))
			}
		}
	}
}

func (s *Service) ApplyToEntity(entity *Entity) {
	s.ServiceBase.ApplyToEntity(entity)
	entity.Service = s.Name
}

func (s *ServiceBase) ApplyToEntity(entity *Entity) {
	if entity.EntityPath == "" {
		entity.EntityPath = s.EntityPath
	}
	if entity.ProtoTypeFile == "" {
		entity.ProtoTypeFile = s.ProtoTypeFile
	}
	if entity.ProtoTypeFileIndent == "" {
		entity.ProtoTypeFileIndent = s.ProtoTypeFileIndent
	}
	if entity.EntityOptPkg == "" {
		entity.EntityOptPkg = s.EntityOptPkg
	}
	if entity.OutputCopierProtoPkgSuffix == "" {
		entity.OutputCopierProtoPkgSuffix = s.OutputCopierProtoPkgSuffix
	}
	entity.Env = mergeEnv(entity.Env, s.Env)
}

func (p *EntityPrefab) ApplyToPrefab(prefab *EntityPrefab) {
	if prefab.Comment.Doc == "" && p.Comment.Doc != "" {
		prefab.Comment.Doc = p.Comment.Doc
	}
	if prefab.Name == "" && p.Name != "" {
		prefab.Name = p.Name
	}
	if prefab.Pkg == "" && p.Pkg != "" {
		prefab.Pkg = p.Pkg
	}
	for _, f := range p.Fields {
		find := false
		for _, v := range prefab.Fields {
			if v.Name == f.Name {
				find = true
				break
			}
		}
		if !find {
			prefab.Fields = append(prefab.Fields, f)
		}
	}
	if prefab.Service == "" && p.Service != "" {
		prefab.Service = p.Service
	}
	if prefab.GrpcSubPkg == "" && p.GrpcSubPkg != "" {
		prefab.GrpcSubPkg = p.GrpcSubPkg
	}
	for _, f := range p.Tmpl {
		find := false
		for _, v := range prefab.Tmpl {
			if v.To == f.To {
				find = true
				break
			}
		}
		if !find {
			prefab.Tmpl = append(prefab.Tmpl, f)
		}
	}
	prefab.Env = mergeEnv(prefab.Env, p.Env)
	if prefab.EntityPath == "" && p.EntityPath != "" {
		prefab.EntityPath = p.EntityPath
	}
	if prefab.ProtoTypeFile == "" && p.ProtoTypeFile != "" {
		prefab.ProtoTypeFile = p.ProtoTypeFile
	}
	if prefab.OptFilePath == "" && p.OptFilePath != "" {
		prefab.OptFilePath = p.OptFilePath
	}
}

func (p *EntityPrefab) ApplyToEntity(entity *Entity) {
	if entity.Comment.Doc == "" && p.Comment.Doc != "" {
		entity.Comment.Doc = p.Comment.Doc
	}
	if entity.Name == "" && p.Name != "" {
		entity.Name = p.Name
	}
	if entity.Pkg == "" && p.Pkg != "" {
		entity.Pkg = p.Pkg
	}
	for _, f := range p.Fields {
		find := false
		for _, v := range entity.Fields {
			if v.Name == f.Name {
				find = true
				break
			}
		}
		if !find {
			entity.Fields = append(entity.Fields, f)
		}
	}
	if entity.Service == "" && p.Service != "" {
		entity.Service = p.Service
	}
	if entity.GrpcSubPkg == "" && p.GrpcSubPkg != "" {
		entity.GrpcSubPkg = p.GrpcSubPkg
	}
	for _, f := range p.Tmpl {
		find := false
		for _, v := range entity.Tmpl {
			if v.To == f.To {
				find = true
				break
			}
		}
		if !find {
			entity.Tmpl = append(entity.Tmpl, f)
		}
	}
	entity.Env = mergeEnv(entity.Env, p.Env)
	if entity.EntityPath == "" && p.EntityPath != "" {
		entity.EntityPath = p.EntityPath
	}
	if entity.ProtoTypeFile == "" && p.ProtoTypeFile != "" {
		entity.ProtoTypeFile = p.ProtoTypeFile
	}
	if entity.OptFilePath == "" && p.OptFilePath != "" {
		entity.OptFilePath = p.OptFilePath
	}
}

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	res := &Config{}
	err = decodeFromYaml(string(content), res)
	if err != nil {
		return nil, err
	}
	res.AfterLoad()
	return res, nil
}

func decodeFromYaml(content string, cfg *Config) error {
	err := yaml.Unmarshal([]byte(content), cfg)
	if err != nil {
		return err
	}
	return nil
}

func mergeEnv(from map[string]map[string]string, to map[string]map[string]string) map[string]map[string]string {
	res := make(map[string]map[string]string)
	for k, v := range to {
		res[k] = map[string]string{}
		for k2, v2 := range v {
			res[k][k2] = v2
		}
	}
	for k, v := range from {
		if _, ok := res[k]; ok {
			for k2, v2 := range v {
				res[k][k2] = v2
			}
		} else {
			res[k] = map[string]string{}
			for k2, v2 := range v {
				res[k][k2] = v2
			}
		}
	}
	return res
}
