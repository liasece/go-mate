package writer

import (
	"reflect"

	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
)

type EntityWriter struct {
	entity reflect.Type
}

func NewEntityWriterByObj(i interface{}) *EntityWriter {
	return &EntityWriter{
		entity: reflect.TypeOf(i),
	}
}

func (w *EntityWriter) GetTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	fs := make([]gocoder.Field, 0)
	for i := 0; i < w.entity.NumField(); i++ {
		ft := w.entity.Field(i)
		fs = append(fs, cde.Field(ft.Name, ft.Type, string(ft.Tag)))
	}
	c.C(cde.Struct(w.entity.Name(), fs...))
	return c
}
