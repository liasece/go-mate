package context

type ConfigTmplContext struct {
	VEntityName  string
	VServiceName string
}

func (c *ConfigTmplContext) EntityName() string {
	return c.VEntityName
}

func (c *ConfigTmplContext) ServiceName() string {
	return c.VServiceName
}
