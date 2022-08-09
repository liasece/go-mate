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
	From  string       `json:"from" yaml:"from"`
	To    string       `json:"to" yaml:"to"`
	Type  TmplItemType `json:"type,omitempty" yaml:"type,omitempty"`
	Merge bool         `json:"merge,omitempty" yaml:"merge,omitempty"`
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
	Comment `json:",inline" yaml:",inline"`
	Base    `json:"base" yaml:"base"`
	Entity  []*Entity `json:"entity" yaml:"entity"`
}

type Entity struct {
	Comment     `json:",inline" yaml:",inline"`
	Name        string         `json:"name" yaml:"name"`
	Pkg         string         `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields      []*EntityField `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service     string         `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg  string         `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	ServiceBase `json:",inline" yaml:",inline"`
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
	for k, service := range c.Service {
		service.Name = k
		for _, tmpl := range service.Tmpl {
			tmpl.AfterLoad()
		}
		service.ServiceBase.AfterLoad()
	}
	for _, entity := range c.Entity {
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
	var data interface{}
	err := yaml.Unmarshal([]byte(content), cfg)
	if err != nil {
		return err
	}
	return decodeFromInterface(data, cfg)
}

func decodeFromInterface(content interface{}, cfg *Config) error {
	switch content := content.(type) {
	case map[string]interface{}:
		return decodeFromMap(content, cfg)
	}
	return nil
}

func decodeFromMap(content map[string]interface{}, cfg *Config) error {
	for k, v := range content {
		switch k {
		case "entity":
			switch v := v.(type) {
			case []interface{}:
				to := make([]*Entity, 0)
				if err := decodeEntitySlice(v, &to); err != nil {
					return err
				}
				cfg.Entity = to
			}
		}
	}
	return nil
}

func decodeEntitySlice(content []interface{}, to *[]*Entity) error {
	for _, v := range content {
		entity := &Entity{}
		if err := decodeEntity(v.(map[string]interface{}), entity); err != nil {
			return err
		}
		*to = append(*to, entity)
	}
	return nil
}

func decodeEntity(content map[string]interface{}, to *Entity) error {
	return nil
}
