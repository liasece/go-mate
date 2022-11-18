package context

import (
	"fmt"

	"github.com/liasece/go-mate/src/gogen/writer/repo"
	"github.com/liasece/go-mate/src/utils"
	"github.com/liasece/gocoder"
)

type EntityStructTmplContext struct {
	utils.TmplUtilsFunc
	w      *repo.RepositoryWriter
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

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) HasFieldsWithTag2(tagReg string, tagReg2 string) bool {
	return len(e.FieldsWithTag(tagReg)) > 0 || len(e.FieldsWithTag(tagReg2)) > 0
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) HasFieldsWithTag3(tagReg string, tagReg2 string, tagReg3 string) bool {
	return len(e.FieldsWithTag(tagReg)) > 0 || len(e.FieldsWithTag(tagReg2)) > 0 || len(e.FieldsWithTag(tagReg3)) > 0
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) HasFieldsWithTagName(tagName string, l1Name string) bool {
	return len(e.FieldsWithTag(fmt.Sprintf("%s:%s", tagName, l1Name))) > 0
}

// tagReg like: `gomate:url` match: `gomate:"foo,url"`
func (e *EntityStructTmplContext) HasFieldsWithTagNameL2(tagName string, l1Name string, l2Name string) bool {
	return len(e.FieldsWithTag(fmt.Sprintf("%s:%s:%s", tagName, l1Name, l2Name))) > 0
}
