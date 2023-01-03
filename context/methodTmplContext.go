package context

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/liasece/gocoder"
)

type MethodTmplContext struct {
	*TmplContext
	method  gocoder.Func
	args    []*ArgTmplContext
	returns []*ArgTmplContext
}

func NewMethodTmplContextList(ctx *TmplContext, methods []gocoder.Func) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewMethodTmplContext(ctx, m))
	}
	return res
}

func NewMethodTmplContext(ctx *TmplContext, method gocoder.Func) *MethodTmplContext {
	return &MethodTmplContext{
		TmplContext: ctx,
		method:      method,
		args:        NewArgTmplContextList(ctx, method.GetArgs()),
		returns:     NewArgTmplContextList(ctx, method.GetReturns()),
	}
}

func (c *MethodTmplContext) Name() string {
	return c.method.GetName()
}

func (c *MethodTmplContext) IsNameReg(nameReg string) bool {
	reg := regexp.MustCompile(nameReg)
	return reg.MatchString(c.Name())
}

func (c *MethodTmplContext) Args() []*ArgTmplContext {
	return c.args
}

func (c *MethodTmplContext) Returns() []*ArgTmplContext {
	return c.returns
}

func (c *MethodTmplContext) GraphqlArgsDefinition() string {
	res := ""
	for _, arg := range c.args {
		if arg.Type().Name() == "error" || arg.Name() == "opUserID" || arg.Type().Name() == "Context" {
			continue
		}
		if arg.Type().IsStruct() && strings.HasSuffix(arg.Type().Name(), "Input") {
			res += arg.Type().FieldsGraphqlDefinition()
			continue
		}
		typeStr := arg.GraphqlType()
		if typeStr == "" {
			continue
		}
		res += fmt.Sprintf("  %s: %s\n", arg.Name(), typeStr)
	}
	return strings.TrimSpace(res)
}

func (c *MethodTmplContext) GraphqlReturnsDefinition() string {
	res := ""
	for _, arg := range c.returns {
		if arg.Type().Name() == "error" || arg.Type().Name() == "Context" {
			continue
		}
		typeStr := arg.GraphqlType()
		if typeStr == "" {
			continue
		}
		res += fmt.Sprintf("  %s: %s\n", arg.Name(), typeStr)
	}
	return strings.TrimSpace(res)
}
