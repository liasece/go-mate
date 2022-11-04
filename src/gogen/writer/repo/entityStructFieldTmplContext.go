package repo

import (
	"reflect"

	"github.com/liasece/go-mate/src/gogen/utils"
	"github.com/liasece/gocoder"
)

type EntityStructFieldTmplContext struct {
	utils.TmplUtilsFunc
	w     *RepositoryWriter
	Field gocoder.Field
}

func (e *EntityStructFieldTmplContext) Name() string {
	return e.Field.GetName()
}

func (e *EntityStructFieldTmplContext) Tag() reflect.StructTag {
	return reflect.StructTag(e.Field.GetTag())
}

func (e *EntityStructFieldTmplContext) Type() *TypeTmplContext {
	return &TypeTmplContext{
		Type: e.Field.GetType(),
	}
}

func (e *EntityStructFieldTmplContext) IsMatchTag(tagReg string) bool {
	return utils.TagMatch(tagReg, string(e.Tag()))
}

func (e *EntityStructFieldTmplContext) IsMatchTag2(tagReg string, tagReg2 string) bool {
	return utils.TagMatch(tagReg, string(e.Tag())) || utils.TagMatch(tagReg2, string(e.Tag()))
}

func (e *EntityStructFieldTmplContext) IsMatchTag3(tagReg string, tagReg2 string, tagReg3 string) bool {
	return utils.TagMatch(tagReg, string(e.Tag())) || utils.TagMatch(tagReg2, string(e.Tag())) || utils.TagMatch(tagReg3, string(e.Tag()))
}

func (e *EntityStructFieldTmplContext) GetMatchTag(tagReg string) string {
	return utils.GetTagMatch(tagReg, string(e.Tag()))
}
