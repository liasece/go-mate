package context

import (
	"github.com/liasece/gocoder"
)

type MethodsTmplContext struct {
	*TmplContext
	methods []*MethodTmplContext
}

func NewMethodsTmplContext(ctx *TmplContext, methods []gocoder.Func) *MethodsTmplContext {
	return &MethodsTmplContext{
		TmplContext: ctx,
		methods:     NewMethodTmplContextList(ctx, methods),
	}
}

func (c *MethodsTmplContext) Methods() []*MethodTmplContext {
	return c.methods
}

func (c *MethodsTmplContext) FindMethods(nameReg string) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0)
	for _, m := range c.methods {
		if m.IsNameReg(nameReg) {
			res = append(res, m)
		}
	}
	return res
}
