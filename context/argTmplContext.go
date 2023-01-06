package context

import (
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
)

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
