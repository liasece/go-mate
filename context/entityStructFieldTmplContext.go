package context

import (
	"reflect"

	"github.com/liasece/go-mate/gogen/writer/repo"
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
)

type EntityStructFieldTmplContext struct {
	*TmplContext
	w     *repo.RepositoryWriter
	Field gocoder.Field
}

func NewEntityStructFieldTmplContextList(ctx *TmplContext, methods []gocoder.Field) []*EntityStructFieldTmplContext {
	res := make([]*EntityStructFieldTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewEntityStructFieldTmplContext(ctx, m))
	}
	return res
}

func NewEntityStructFieldTmplContext(ctx *TmplContext, typ gocoder.Field) *EntityStructFieldTmplContext {
	return &EntityStructFieldTmplContext{
		TmplContext: ctx,
		Field:       typ,
	}
}

func (e *EntityStructFieldTmplContext) Name() string {
	return e.Field.GetName()
}

func (e *EntityStructFieldTmplContext) Tag() reflect.StructTag {
	return reflect.StructTag(e.Field.GetTag())
}

func (e *EntityStructFieldTmplContext) Type() *TypeTmplContext {
	return &TypeTmplContext{
		TmplContext: nil,
		Type:        e.Field.GetType(),
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

func (e *EntityStructFieldTmplContext) GraphqlDefinition() string {
	return utils.ToLowerCamelCase(e.Name()) + ": " + utils.GraphqlStyle(e.Name(), e.Type().Name())
}
