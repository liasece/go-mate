package context

type RawTmplContext struct {
	*TmplContext
	raw map[string]interface{}
}

func NewRawTmplContext(ctx *TmplContext, raw map[string]interface{}) *RawTmplContext {
	return &RawTmplContext{
		TmplContext: ctx,
		raw:         raw,
	}
}

func (e *RawTmplContext) Raw() map[string]interface{} {
	return e.raw
}
