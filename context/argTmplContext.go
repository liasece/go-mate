package context

import (
	"fmt"
	"strings"

	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
)

var _ IField = (*ArgTmplContext)(nil)

type ArgTmplContext struct {
	*TmplContext
	gocoder.Arg
	typ *TypeTmplContext
}

func NewArgTmplContextList(ctx *TmplContext, methods []gocoder.Arg) []*ArgTmplContext {
	res := make([]*ArgTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewArgTmplContext(ctx, m))
	}
	return res
}

func NewArgTmplContext(ctx *TmplContext, arg gocoder.Arg) *ArgTmplContext {
	return &ArgTmplContext{
		TmplContext: ctx,
		Arg:         arg,
		typ:         NewTypeTmplContext(ctx, arg.GetType()),
	}
}

func (e *ArgTmplContext) GetTmplContext() *TmplContext {
	return e.TmplContext
}

func (e *ArgTmplContext) Doc() string {
	resList := make([]string, 0)
	for _, note := range e.Arg.Notes() {
		resList = append(resList, note.GetContent())
	}
	return strings.Join(resList, "\n")
}

func (e *ArgTmplContext) Name() string {
	return e.Arg.GetName()
}

func (e *ArgTmplContext) Type() *TypeTmplContext {
	return e.typ
}

func (e *ArgTmplContext) GraphqlType() string {
	return utils.GraphqlStyle(e.Name(), e.typ.Name())
}

func (e *ArgTmplContext) ProtoBuffType() string {
	return utils.ProtoBuffTypeStyle(e.Name(), e.typ.Name())
}

// return like `foo: [String!]`
func (e *ArgTmplContext) GraphqlDefinition() string {
	typeStr := e.GraphqlType()
	if typeStr == "" {
		return ""
	}
	return utils.ToLowerCamelCase(e.Name()) + ": " + typeStr
}

// return like `repeated string foo = 1;`
func (e *ArgTmplContext) ProtoBuffDefinition(argIndex int) string {
	typeStr := e.ProtoBuffType()
	if typeStr == "" {
		return ""
	}
	name := e.Name()
	if name == "" {
		name = fmt.Sprintf("arg%d", argIndex)
	}
	return fmt.Sprintf("%s %s = %d;", typeStr, utils.SnakeString(name), argIndex)
}
