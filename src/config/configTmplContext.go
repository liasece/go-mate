package config

import "github.com/liasece/go-mate/src/gogen/utils"

type ConfigTmplContext struct {
	utils.TmplUtilsFunc
	VEntityName  string
	VServiceName string
}

func (c *ConfigTmplContext) EntityName() string {
	return c.VEntityName
}

func (c *ConfigTmplContext) ServiceName() string {
	return c.VServiceName
}
