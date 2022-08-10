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

func (e *EntityStructFieldTmplContext) Type() string {
	return e.Field.GetType().String()
}

func (e *EntityStructFieldTmplContext) TypeIsNumber() bool {
	return e.Field.GetType().Kind() >= reflect.Int && e.Field.GetType().Kind() <= reflect.Float64
}
