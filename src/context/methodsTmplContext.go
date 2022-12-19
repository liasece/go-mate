package context

import (
	"github.com/liasece/gocoder"
)

type MethodsTmplContext struct {
	*TmplContext
	Methods []gocoder.Func
}

func NewMethodsTmplContext(ctx *TmplContext, methods []gocoder.Func) *MethodsTmplContext {
	return &MethodsTmplContext{
		TmplContext: ctx,
		Methods:     methods,
	}
}
