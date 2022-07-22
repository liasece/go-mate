package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Comment struct {
	Doc string `json:"doc" yaml:"doc"`
}

type Config struct {
	Comment `json:",inline" yaml:",inline"`
	Entity  []*Entity `json:"entity" yaml:"entity"`
}

type Entity struct {
	Comment `json:",inline" yaml:",inline"`
	Name    string         `json:"name" yaml:"name"`
	Pkg     string         `json:"pkg" yaml:"pkg"`
	Fields  []*EntityField `json:"fields" yaml:"fields"`
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

func LoadConfig(path string) (*Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var res Config
	err = decodeFromYaml(string(content), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
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
