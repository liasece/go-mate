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

// check this method doc is match the reg, docReg like `@ext\s+@graphql\s+@mutation.*`, doc like `@ext @graphql @mutation`
func (c *MethodTmplContext) IsDocReg(docReg string) bool {
	reg := regexp.MustCompile(docReg)
	for _, note := range c.method.Notes() {
		if reg.MatchString(note.GetContent()) {
			return true
		}
	}
	return false
}

// docReg like `@description\s+(.*)` group like 1, doc like `@description xxx`, return `xxx`
func (c *MethodTmplContext) GetDocByReg(docReg string, group int) string {
	reg := regexp.MustCompile(docReg)
	for _, note := range c.method.Notes() {
		if ss := reg.FindStringSubmatch(note.GetContent()); len(ss) > 0 {
			return ss[group]
		}
	}
	return ""
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
		argType := NewTypeTmplContext(c.TmplContext, arg.Type().UnPtr())
		if argType.IsStruct() && strings.HasSuffix(argType.Name(), "Input") {
			res += argType.FieldsGraphqlDefinition()
			continue
		}
		typeStr := arg.GraphqlType()
		if typeStr == "" {
			continue
		}
		{
			// add doc
			if doc := c.GetDocByReg(`@return\s+`+arg.Name()+`\s+(.*)`, 1); doc != "" {
				res += fmt.Sprintf("  \"\"\"\n%s\n\"\"\"\n", doc)
			}
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
		{
			// add doc
			if doc := c.GetDocByReg(`@return\s+`+arg.Name()+`\s+(.*)`, 1); doc != "" {
				res += fmt.Sprintf("  \"\"\"\n%s\n\"\"\"\n", doc)
			}
		}
		res += fmt.Sprintf("  %s: %s\n", arg.Name(), typeStr)
	}
	return strings.TrimSpace(res)
}
