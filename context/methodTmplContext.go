package context

import (
	"regexp"
	"strings"

	"github.com/liasece/gocoder"
)

type MethodTmplContext struct {
	*TmplContext
	method  gocoder.Func
	args    []*ArgTmplContext
	returns []*ArgTmplContext

	argFieldsTmpl    *FieldsTmplContext
	returnFieldsTmpl *FieldsTmplContext
}

func NewMethodTmplContextList(ctx *TmplContext, methods []gocoder.Func) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewMethodTmplContext(ctx, m))
	}
	return res
}

func NewMethodTmplContext(ctx *TmplContext, method gocoder.Func) *MethodTmplContext {
	args := NewArgTmplContextList(ctx, method.GetArgs())
	returns := NewArgTmplContextList(ctx, method.GetReturns())
	res := &MethodTmplContext{
		TmplContext:      ctx,
		method:           method,
		args:             args,
		returns:          returns,
		argFieldsTmpl:    nil,
		returnFieldsTmpl: nil,
	}
	{
		argFields := make([]IField, 0, len(args))
		for _, v := range args {
			argFields = append(argFields, v)
		}
		argFieldsTmpl := NewFieldsTmplContext(ctx, argFields)
		argFieldsTmpl.docReader = res.GetParamStdDoc

		res.argFieldsTmpl = argFieldsTmpl
	}
	{
		returnFields := make([]IField, 0, len(returns))
		for _, v := range returns {
			returnFields = append(returnFields, v)
		}
		returnFieldsTmpl := NewFieldsTmplContext(ctx, returnFields)
		returnFieldsTmpl.docReader = res.GetReturnStdDoc

		res.returnFieldsTmpl = returnFieldsTmpl
	}
	return res
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

// GetStdDoc like `@<docGroup>\s+<fieldName>\s+(.*)` group 1, doc like `@param foo xxx`, docGroup = param, fieldName == foo, return `xxx`
func (c *MethodTmplContext) GetStdDoc(docGroup string, fieldName string) string {
	return c.GetDocByReg(`@`+regexp.QuoteMeta(docGroup)+`\s+`+regexp.QuoteMeta(fieldName)+`\s+(.*)`, 1)
}

// alias GetStdDoc("param", fieldName)
func (c *MethodTmplContext) GetParamStdDoc(fieldName string) string {
	return c.GetStdDoc("param", fieldName)
}

// alias GetStdDoc("return", fieldName)
func (c *MethodTmplContext) GetReturnStdDoc(fieldName string) string {
	return c.GetStdDoc("return", fieldName)
}

func (c *MethodTmplContext) GraphqlArgsDefinition() string {
	return c.argFieldsTmpl.GraphqlDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Name() != "opUserID" && i.Type().Name() != "Context"
	})
}

func (c *MethodTmplContext) GraphqlReturnsDefinition() string {
	return c.returnFieldsTmpl.GraphqlDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Type().Name() != "Context"
	})
}

func (c *MethodTmplContext) ProtoBuffArgsDefinition() string {
	return c.argFieldsTmpl.ProtoBuffDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Type().Name() != "Context"
	})
}

func (c *MethodTmplContext) ProtoBuffReturnsDefinition() string {
	return c.returnFieldsTmpl.ProtoBuffDefinitionFilterFunc(func(i IField) bool {
		return i.Name() != "error"
	})
}

func (c *MethodTmplContext) GRPCCallGoArgsDefinition(reqValueName string) string {
	return c.argFieldsTmpl.GRPCCallGoDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "Context"
	}, reqValueName)
}

func (c *MethodTmplContext) GoArgsDefinition() string {
	return c.argFieldsTmpl.GoDefinition()
}

func (c *MethodTmplContext) GraphqlGoArgsDefinition() string {
	return c.argFieldsTmpl.GraphqlGoDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Name() != "opUserID"
	})
}

func (c *MethodTmplContext) CallGRPCArgsDefinition() string {
	return c.argFieldsTmpl.CallGRPCDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "Context"
	})
}

func (c *MethodTmplContext) CallGoReturnsDefinition() string {
	return c.returnFieldsTmpl.CallGoDefinitionFilterFunc(func(i IField) bool {
		return true
	})
}

func (c *MethodTmplContext) GRPCCallGoReturnsResponseDefinition() string {
	return c.returnFieldsTmpl.GRPCCallGoResponseDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error"
	})
}

func (c *MethodTmplContext) IsInternal() bool {
	// internal method name start with lower case
	return strings.ToLower(c.Name()[0:1]) == c.Name()[0:1]
}
