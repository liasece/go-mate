package context

import (
	"github.com/liasece/go-mate/src/config"
)

type ITmplContext interface {
	Terminate() bool
	GetTerminate() bool
}

type BaseTmplContext struct {
	terminate bool
}

func (e *BaseTmplContext) Terminate() bool {
	e.terminate = true
	return e.terminate
}

func (e *BaseTmplContext) GetTerminate() bool {
	return e.terminate
}

type TmplContext struct {
	BaseTmplContext
	tmpl *config.TmplItem
}

func NewTmplContext(tmpl *config.TmplItem) *TmplContext {
	return &TmplContext{
		BaseTmplContext: BaseTmplContext{},
		tmpl:            tmpl,
	}
}

func (e *TmplContext) TmplFrom() string {
	return e.tmpl.From
}

func (e *TmplContext) TmplTo() string {
	return e.tmpl.To
}

func (e *TmplContext) TmplType() string {
	return string(e.tmpl.Type)
}

func (e *TmplContext) TmplMerge() bool {
	return e.tmpl.Merge
}

func (e *TmplContext) TmplOnlyCreate() bool {
	return e.tmpl.OnlyCreate
}

func (e *TmplContext) TmplOptEq(key string, value string) bool {
	return e.TmplOpt(key) == value
}

func (e *TmplContext) TmplOpt(key string) string {
	if res, ok := e.tmpl.Opt[key]; ok {
		return res
	}
	return ""
}
