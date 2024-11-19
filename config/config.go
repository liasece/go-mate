package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/liasece/gocoder"
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
	From       string            `json:"from" yaml:"from"`
	To         string            `json:"to" yaml:"to"`
	Type       TmplItemType      `json:"type,omitempty" yaml:"type,omitempty"`
	Merge      bool              `json:"merge,omitempty" yaml:"merge,omitempty"`
	OnlyCreate bool              `json:"onlyCreate,omitempty" yaml:"onlyCreate,omitempty"`
	Opt        map[string]string `json:"opt,omitempty" yaml:"opt,omitempty"`
}

type ServiceBase struct {
	EntityPath                 string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	EntityKind                 string                       `json:"entityKind,omitempty" yaml:"entityKind,omitempty"`
	EntityRealName             string                       `json:"entityRealName,omitempty" yaml:"entityRealName,omitempty"`
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
	Service []*Service `json:"service" yaml:"service"`
}

type Comment struct {
	Doc string `json:"doc,omitempty" yaml:"doc,omitempty"`
}

type Config struct {
	Comment               `json:",inline" yaml:",inline"`
	Base                  `json:"base" yaml:"base"`
	Entity                []*Entity                    `json:"entity,omitempty" yaml:"entity,omitempty"`
	EntityPrefab          []*EntityPrefab              `json:"entityPrefab,omitempty" yaml:"entityPrefab,omitempty"`
	BuildEntityWithPrefab yaml.MapSlice                `json:"buildEntityWithPrefab,omitempty" yaml:"buildEntityWithPrefab,omitempty"`
	Env                   map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	LogLevel              string                       `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	ImportGoCodePath      []string                     `json:"importGoCodePath,omitempty" yaml:"importGoCodePath,omitempty"`
	Raw                   *Raw                         `json:"raw,omitempty" yaml:"raw,omitempty"`
}

type Raw struct {
	Config string      `json:"config" yaml:"config"`
	Tmpl   []*TmplItem `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
}

type EntityPrefab struct {
	Comment        `json:",inline" yaml:",inline"`
	Name           string                       `json:"name" yaml:"name"`
	Pkg            string                       `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields         []*EntityField               `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service        string                       `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg     string                       `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	Prefab         []string                     `json:"prefab,omitempty" yaml:"prefab,omitempty"`
	Tmpl           []*TmplItem                  `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
	Env            map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	EntityPath     string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	EntityKind     string                       `json:"entityKind,omitempty" yaml:"entityKind,omitempty"`
	EntityRealName string                       `json:"entityRealName,omitempty" yaml:"entityRealName,omitempty"`
	RepeatByPrefab []string                     `json:"repeatByPrefab,omitempty" yaml:"repeatByPrefab,omitempty"`
	ProtoTypeFile  string                       `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	NoSelector     *bool                        `json:"noSelector,omitempty" yaml:"noSelector,omitempty"`
}

type Entity struct {
	Comment                    `json:",inline" yaml:",inline"`
	Name                       string                       `json:"name" yaml:"name"`
	CodeName                   string                       `json:"-" yaml:"-"`
	Pkg                        string                       `json:"pkg,omitempty" yaml:"pkg,omitempty"`
	Fields                     []*EntityField               `json:"fields,omitempty" yaml:"fields,omitempty"`
	Service                    string                       `json:"service,omitempty" yaml:"service,omitempty"`
	GrpcSubPkg                 string                       `json:"grpcSubPkg,omitempty" yaml:"grpcSubPkg,omitempty"`
	Prefab                     []string                     `json:"prefab,omitempty" yaml:"prefab,omitempty"`
	Tmpl                       []*TmplItem                  `json:"tmpl,omitempty" yaml:"tmpl,omitempty"`
	Env                        map[string]map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	ProtoTypeFileIndent        string                       `json:"protoTypeFileIndent,omitempty" yaml:"protoTypeFileIndent,omitempty"`
	EntityPath                 string                       `json:"entityPath,omitempty" yaml:"entityPath,omitempty"`
	EntityKind                 string                       `json:"entityKind,omitempty" yaml:"entityKind,omitempty"`
	EntityRealName             string                       `json:"entityRealName,omitempty" yaml:"entityRealName,omitempty"`
	DecodedEntityPath          string                       `json:"-" yaml:"-"`
	ProtoTypeFile              string                       `json:"protoTypeFile,omitempty" yaml:"protoTypeFile,omitempty"`
	EntityOptPkg               string                       `json:"entityOptPkg,omitempty" yaml:"entityOptPkg,omitempty"`
	OutputCopierProtoPkgSuffix string                       `json:"outputCopierProtoPkgSuffix,omitempty" yaml:"outputCopierProtoPkgSuffix,omitempty"`
	NoSelector                 *bool                        `json:"noSelector,omitempty" yaml:"noSelector,omitempty"`

	CodeType       gocoder.Codable `json:"-" yaml:"-"`
	ConfigFilePath string          `json:"-" yaml:"-"`
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

// get entity config env
func (c *Entity) GetEnv(k1, k2 string) string {
	if c.Env == nil {
		return ""
	}
	if v, ok := c.Env[k1]; ok {
		if v2, ok := v[k2]; ok {
			return v2
		}
	}
	return ""
}

func (c *TmplItem) AfterLoad() {
	if c.Type == "" {
		switch {
		case strings.HasSuffix(c.To, ".go"):
			c.Type = TmplItemTypeGo
		case strings.HasSuffix(c.To, ".proto"):
			c.Type = TmplItemTypeProto
		case strings.HasSuffix(c.To, ".graphql"):
			c.Type = TmplItemTypeGraphQL
		default:
			log.Fatal("TmplItem AfterLoad unknown tmpl type", log.Any("tmpl", c))
		}
	}
}

func (s *ServiceBase) AfterLoad() {
	getIndent := func() string {
		if s.ProtoTypeFileIndent == "" {
			return "\t"
		}
		if strings.HasPrefix(s.ProtoTypeFileIndent, "$") {
			switch s.ProtoTypeFileIndent[1:] {
			case "tab":
				return "\t"
			default:
				a, _ := strconv.Atoi(s.ProtoTypeFileIndent[1:])
				return strings.Repeat(" ", a)
			}
		}
		return s.ProtoTypeFileIndent
	}
	s.ProtoTypeFileIndent = getIndent()
}

func (c *Config) getPrefab(name string) *EntityPrefab {
	for _, p := range c.EntityPrefab {
		if p.Name == name {
			return p
		}
	}
	return nil
}

func interfaceToStringList(v interface{}) []string {
	if v == nil {
		return nil
	}
	switch v := v.(type) {
	case []string:
		return v
	case []interface{}:
		var ret []string
		for _, v := range v {
			ret = append(ret, v.(string))
		}
		return ret
	}
	return nil
}

func (c *Config) AfterLoad() {
	// build prefab
	for _, prefab := range c.EntityPrefab {
		prefab.Env = mergeEnv(prefab.Env, c.Env)
	}

	for ptr, prefab := range c.EntityPrefab {
		for _, innerPrefab := range prefab.Prefab {
			var findEntityPrefab *EntityPrefab
			for index, v := range c.EntityPrefab {
				if v.Name == innerPrefab {
					findEntityPrefab = v
					if index > ptr {
						log.Error("Prefab dependency order is wrong, "+prefab.Name+" tries to reference uninitialized "+v.Name, log.Any("prefab", prefab), log.Any("innerPrefab", innerPrefab))
					}
					break
				}
			}
			if findEntityPrefab != nil {
				findEntityPrefab.ApplyToPrefab(prefab)
			} else {
				log.Fatal("Config AfterLoad entity prefab not found", log.Any("innerPrefab", innerPrefab))
			}
		}
		for _, tmpl := range prefab.Tmpl {
			tmpl.AfterLoad()
		}
	}

	{
		// apply repeat by prefab
		for _, originBuildEntityWithPrefabV := range c.BuildEntityWithPrefab {
			originPrefab := c.getPrefab(originBuildEntityWithPrefabV.Key.(string))
			for _, repeatByPrefabName := range originPrefab.RepeatByPrefab {
				find := false
				for i, buildEntityWithPrefabV := range c.BuildEntityWithPrefab {
					if buildEntityWithPrefabV.Key.(string) == repeatByPrefabName {
						c.BuildEntityWithPrefab[i] = yaml.MapItem{
							Key:   buildEntityWithPrefabV.Key,
							Value: append(interfaceToStringList(buildEntityWithPrefabV.Value), interfaceToStringList(originBuildEntityWithPrefabV.Value)...),
						}
						find = true
					}
				}
				if !find {
					c.BuildEntityWithPrefab = append(c.BuildEntityWithPrefab, yaml.MapItem{
						Key:   repeatByPrefabName,
						Value: originBuildEntityWithPrefabV.Value,
					})
				}
			}
		}
	}

	for _, v := range c.BuildEntityWithPrefab {
		prefabName, entityNameList := v.Key.(string), interfaceToStringList(v.Value)
		prefab := c.getPrefab(prefabName)
		for _, entityName := range entityNameList {
			find := false
			for _, v := range c.Entity {
				if v.Name == entityName && v.EntityRealName == prefab.EntityRealName && v.EntityKind == prefab.EntityKind {
					find = true
					break
				}
			}
			if !find {
				c.Entity = append(c.Entity, &Entity{
					Name:           entityName,
					EntityKind:     prefab.EntityKind,
					EntityRealName: prefab.EntityRealName,
					Prefab:         []string{prefabName},

					// it will be override by prefab list
					Comment:                    Comment{Doc: ""},
					Pkg:                        "",
					Fields:                     nil,
					Service:                    "",
					GrpcSubPkg:                 "",
					Tmpl:                       nil,
					Env:                        nil,
					ProtoTypeFileIndent:        "",
					EntityPath:                 "",
					DecodedEntityPath:          "",
					ProtoTypeFile:              "",
					EntityOptPkg:               "",
					OutputCopierProtoPkgSuffix: "",
					NoSelector:                 nil,
					CodeType:                   nil,
				})
			} else {
				log.Fatal("Config AfterLoad build prefab target entity already exists", log.Any("entityName", entityName))
			}
		}
	}

	for _, service := range c.Service {
		service.ServiceBase.AfterLoad()
		service.Env = mergeEnv(service.Env, c.Env)
	}
	for _, entity := range c.Entity {
		for _, prefab := range entity.Prefab {
			var findEntityPrefab *EntityPrefab
			for _, v := range c.EntityPrefab {
				if v.Name == prefab {
					findEntityPrefab = v
					break
				}
			}
			if findEntityPrefab != nil {
				findEntityPrefab.ApplyToEntity(entity)
			} else {
				log.Fatal("Config AfterLoad entity prefab not found", log.Any("entity", entity))
			}
		}
		if entity.Service != "" {
			var findService *Service
			for _, v := range c.Service {
				if v.Name == entity.Service {
					findService = v
					break
				}
			}
			if findService != nil {
				findService.ApplyToEntity(entity)
			} else {
				log.Fatal("Config AfterLoad service not found", log.Any("entity", entity))
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
	if entity.EntityKind == "" {
		entity.EntityKind = s.EntityKind
	}
	if entity.EntityRealName == "" {
		entity.EntityRealName = s.EntityRealName
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

func mapEq(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

func (p *EntityPrefab) ApplyToPrefab(prefab *EntityPrefab) {
	if prefab.Comment.Doc == "" && p.Comment.Doc != "" {
		prefab.Comment.Doc = p.Comment.Doc
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
			if v.To == f.To && v.From == f.From && mapEq(v.Opt, f.Opt) {
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
	if prefab.EntityKind == "" && p.EntityKind != "" {
		prefab.EntityKind = p.EntityKind
	}
	if prefab.EntityRealName == "" && p.EntityRealName != "" {
		prefab.EntityRealName = p.EntityRealName
	}
	for _, f := range p.RepeatByPrefab {
		find := false
		for _, v := range prefab.RepeatByPrefab {
			if v == f {
				find = true
				break
			}
		}
		if !find {
			prefab.RepeatByPrefab = append(prefab.RepeatByPrefab, f)
		}
	}
	if prefab.ProtoTypeFile == "" && p.ProtoTypeFile != "" {
		prefab.ProtoTypeFile = p.ProtoTypeFile
	}
	if prefab.NoSelector == nil && p.NoSelector != nil {
		prefab.NoSelector = p.NoSelector
	}
}

func (p *EntityPrefab) ApplyToEntity(entity *Entity) {
	if entity.Comment.Doc == "" && p.Comment.Doc != "" {
		entity.Comment.Doc = p.Comment.Doc
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
			if v.To == f.To && v.From == f.From && mapEq(v.Opt, f.Opt) {
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
	if entity.EntityKind == "" && p.EntityKind != "" {
		entity.EntityKind = p.EntityKind
	}
	if entity.EntityRealName == "" && p.EntityRealName != "" {
		entity.EntityRealName = p.EntityRealName
	}
	if entity.ProtoTypeFile == "" && p.ProtoTypeFile != "" {
		entity.ProtoTypeFile = p.ProtoTypeFile
	}
	if entity.NoSelector == nil && p.NoSelector != nil {
		entity.NoSelector = p.NoSelector
	}
}

func LoadConfig(path string) (*Config, error) {
	var res Config
	err := LoadConfigTo(path, &res)
	if err != nil {
		return nil, err
	}
	res.AfterLoad()
	return &res, nil
}

func LoadConfigTo(path string, res interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = decodeFromYaml(string(content), res)
	if err != nil {
		return err
	}
	return nil
}

func decodeFromYaml(content string, cfg interface{}) error {
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
