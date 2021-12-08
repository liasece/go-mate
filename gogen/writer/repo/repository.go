package repo

import (
	"fmt"
	"reflect"

	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"go.mongodb.org/mongo-driver/bson"
)

type RepositoryWriter struct {
	entity     reflect.Type
	entityName string
}

func NewRepositoryWriterByObj(i interface{}) *RepositoryWriter {
	t := reflect.TypeOf(i)
	return &RepositoryWriter{
		entity:     t,
		entityName: t.Name(),
	}
}

func NewRepositoryWriterByType(t reflect.Type, name string) *RepositoryWriter {
	return &RepositoryWriter{
		entity:     t,
		entityName: name,
	}
}

func (w *RepositoryWriter) GetFilterTypeStructCodeStruct() gocoder.Struct {
	res, _ := w.GetFilterTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetUpdaterTypeStructCodeStruct() gocoder.Struct {
	res, _ := w.GetUpdaterTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetFilterTypeStructCode() (gocoder.Struct, []*FieldFilterField) {
	mfs := make([]*FieldFilterField, 0)
	mfs = append(mfs, newFieldFilterField(reflect.StructField{}, "NoCount", cde.Type(false)))
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldFilterFields(w.entity.Field(i))...)
	}
	strT := cde.Struct(fmt.Sprintf("%sFilter", w.entityName), fieldFilterFieldsToGocoder(mfs)...)
	return strT, mfs
}

func (w *RepositoryWriter) GetFilterTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	strT, mfs := w.GetFilterTypeStructCode()
	c.C(strT)
	c.C(getFieldFilterMethodToBSON(strT, mfs))

	return c
}

func (w *RepositoryWriter) GetUpdaterTypeStructCode() (gocoder.Struct, []*FieldUpdaterField) {
	mfs := make([]*FieldUpdaterField, 0)
	mfs = append(mfs, newFieldUpdaterField(reflect.StructField{}, "JustDelete", cde.Type(true)))
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldUpdaterFields(w.entity.Field(i))...)
	}

	strT := cde.Struct(fmt.Sprintf("%sUpdater", w.entityName), fieldUpdaterFieldsToGocoder(mfs)...)
	return strT, mfs
}

func (w *RepositoryWriter) GetUpdaterTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	strT, mfs := w.GetUpdaterTypeStructCode()
	c.C(strT)
	c.C(getFieldUpdaterMethodToBSON(strT, mfs))

	return c
}

func (w *RepositoryWriter) getQueryCode(receiver gocoder.Receiver, filter gocoder.Struct) gocoder.Code {
	c := gocoder.NewCode()
	resV := cde.Value("res", cde.Type(w.entity).Slice())
	countV := cde.Value("count", cde.Type(0))
	errV := cde.Value("err", cde.TypeError())
	queryArgV := cde.Arg("query", filter.GetType())
	ctxV := cde.Arg("ctx", cde.Type("context.Context"))
	method := cde.Method("Query", receiver, []gocoder.Arg{
		ctxV,
		queryArgV,
	}, []gocoder.Type{
		resV.Type(),
		countV.Type(),
		errV.Type(),
	})
	method.C(
		resV.AutoSet(cde.Make(resV, 0)),
		countV.AutoSet(0),
	)
	bsonQueryV := cde.Value("bsonQuery", cde.Type(bson.M{}))
	method.C(
		bsonQueryV.AutoSet(queryArgV.GetValue().Method("ToBSON").Call()),
		// cde.Values(countV, errV).AutoSet(receiver.GetValue().Field("Mongo").Method("Count").Call()),
		cde.Return(resV, countV, cde.Value("nil", nil)),
	)
	c.C(method)
	return c
}

func (w *RepositoryWriter) GetEntityRepositoryCode(filter gocoder.Struct, updater gocoder.Struct) gocoder.Code {
	c := gocoder.NewCode()

	mfs := make([]gocoder.Field, 0)
	mfs = append(mfs)

	strT := cde.Struct(fmt.Sprintf("%sRepository", w.entityName), mfs...)
	c.C(strT)
	receiver := cde.Receiver("r", strT.GetType().TackPtr())
	c.C(w.getQueryCode(receiver, filter))

	return c
}

func (w *RepositoryWriter) GetEntityRepositoryStructCode() gocoder.Code {
	return w.GetEntityRepositoryCode(w.GetFilterTypeStructCodeStruct(), w.GetUpdaterTypeStructCodeStruct())
}
