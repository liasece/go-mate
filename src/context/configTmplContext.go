package context

import "github.com/liasece/go-mate/src/utils"

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
