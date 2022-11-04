package repo

import (
	"github.com/liasece/go-mate/src/gogen/utils"
	"github.com/liasece/gocoder"
)

type EntityStructTmplContext struct {
	utils.TmplUtilsFunc
	w      *RepositoryWriter
	Struct gocoder.Struct
}

func (e *EntityStructTmplContext) Fields() []*EntityStructFieldTmplContext {
	fields := make([]*EntityStructFieldTmplContext, 0)
	for _, field := range e.Struct.GetFields() {
		fields = append(fields, &EntityStructFieldTmplContext{
			w:     e.w,
			Field: field,
		})
	}
	return fields
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) FieldsWithTag(tagReg string) []*EntityStructFieldTmplContext {
	fields := make([]*EntityStructFieldTmplContext, 0)
	for _, field := range e.Struct.GetFields() {
		if !utils.TagMatch(tagReg, field.GetTag()) {
			continue
		}
		fields = append(fields, &EntityStructFieldTmplContext{
			w:     e.w,
			Field: field,
		})
	}
	return fields
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) HasFieldsWithTag(tagReg string) bool {
	return len(e.FieldsWithTag(tagReg)) > 0
}
