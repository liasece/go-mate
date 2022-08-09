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
	TmplItemTypeGo    TmplItemType = "go"
	TmplItemTypeProto TmplItemType = "proto"
)

type TmplItem struct {
	From       string       `json:"from" yaml:"from"`
	To         string       `json:"to" yaml:"to"`
	Type       TmplItemType `json:"type,omitempty" yaml:"type,omitempty"`
	Merge      bool         `json:"merge,omitempty" yaml:"merge,omitempty"`
	OnlyCreate bool         `json:"onlyCreate,omitempty" yaml:"onlyCreate,omitempty"`
}

type ServiceBase struct {
	EntityPath                 string      `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	ProtoTypeFile              string      `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	ProtoTypeFileIndent        string      `json:"protoTypeFileIndent,omitempty" yaml:"protoTypeFileIndent,omitempty"`
	CopierFile                 string      `json:"copierFile,omitempty" yaml:"copierFile,omitempty"`
	Tmpl                       []*TmplItem `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
	EntityOptPkg               string      `json:"entityOptPkg,omitempty" yaml:"entityOptPkg,omitempty"`
	OutputCopierProtoPkgSuffix string      `json:"outputCopierProtoPkgSuffix,omitempty" yaml:"outputCopierProtoPkgSuffix,omitempty"`
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
	Entity                []*Entity                `json:"entity,omitempty" yaml:"entity,omitempty"`
	EntityPrefab          map[string]*EntityPrefab `json:"entityPrefab,omitempty" yaml:"entityPrefab,omitempty"`
	BuildEntityWithPrefab map[string][]string      `json:"buildEntityWithPrefab,omitempty" yaml:"buildEntityWithPrefab,omitempty"`
}

type EntityPrefab struct {
	Comment    `json:",inline" yaml:",inline"`
	Name       string         `json:"name" yaml:"name"`
	Pkg        string         `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields     []*EntityField `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service    string         `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg string         `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
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
}

type Entity struct {
	Comment     `json:",inline" yaml:",inline"`
	Name        string         `json:"name" yaml:"name"`
	Pkg         string         `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields      []*EntityField `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service     string         `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg  string         `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	ServiceBase `json:",inline" yaml:",inline"`
	Prefab      string `json:"prefab,omitempty" yaml:"prefab,omitempty"`
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
					Prefab: prefabName,
				})
			} else {
				log.L(nil).Fatal("Config AfterLoad build prefab target entity already exists", log.Any("entityName", entityName))
			}
		}
	}

	for k, service := range c.Service {
		service.Name = k
		for _, tmpl := range service.Tmpl {
			tmpl.AfterLoad()
		}
		service.ServiceBase.AfterLoad()
	}
	for _, entity := range c.Entity {
		if entity.Prefab != "" {
			if prefab, ok := c.EntityPrefab[entity.Prefab]; ok {
				prefab.ApplyToEntity(entity)
			} else {
				log.L(nil).Fatal("Config AfterLoad entity prefab not found", log.Any("entity", entity))
			}
		}
		if entity.Service != "" {
			if service, ok := c.Service[entity.Service]; ok {
				entity.ServiceBase = service.ServiceBase
			} else {
				log.L(nil).Fatal("Config AfterLoad service not found", log.Any("entity", entity))
			}
		}
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
