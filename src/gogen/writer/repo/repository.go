package repo

import (
	"fmt"

	"github.com/liasece/go-mate/src/config"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/gocoder/cdt"
	"go.mongodb.org/mongo-driver/bson"
)

type RepositoryWriter struct {
	entity      gocoder.Type
	entityName  string
	entityPkg   string
	serviceName string

	Filter   gocoder.Struct
	Updater  gocoder.Struct
	Sorter   gocoder.Struct
	Selector gocoder.Struct

	OutTypeSuffix     string
	OutFilterSuffix   string
	OutUpdaterSuffix  string
	OutSorterSuffix   string
	OutSelectorSuffix string

	EntityCfg *config.Entity
}

func NewRepositoryWriterByObj(i interface{}) *RepositoryWriter {
	t := cde.Type(i)
	return &RepositoryWriter{
		entity:     t,
		entityName: t.Name(),
	}
}

func NewRepositoryWriterByType(t gocoder.Type, name string, pkg string, serviceName string, outFilterSuffix string, outUpdaterSuffix string, outSorterSuffix string, outSelectorSuffix string, outTypeSuffix string) *RepositoryWriter {
	return &RepositoryWriter{
		entity:            t,
		entityName:        name,
		entityPkg:         pkg,
		serviceName:       serviceName,
		OutTypeSuffix:     outTypeSuffix,
		OutFilterSuffix:   outFilterSuffix,
		OutUpdaterSuffix:  outUpdaterSuffix,
		OutSorterSuffix:   outSorterSuffix,
		OutSelectorSuffix: outSelectorSuffix,
	}
}

func (w *RepositoryWriter) GetFilterTypeStructCodeStruct() gocoder.Struct {
	if w.Filter != nil {
		return w.Filter
	}
	res, _ := w.GetFilterTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetUpdaterTypeStructCodeStruct() gocoder.Struct {
	if w.Updater != nil {
		return w.Updater
	}
	res, _ := w.GetUpdaterTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetSorterTypeStructCodeStruct() gocoder.Struct {
	if w.Sorter != nil {
		return w.Sorter
	}
	res, _ := w.GetSorterTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetSelectorTypeStructCodeStruct() gocoder.Struct {
	if w.Selector != nil {
		return w.Selector
	}
	res, _ := w.GetSelectorTypeStructCode()
	return res
}

func (w *RepositoryWriter) GetFilterTypeStructCode() (gocoder.Struct, []*FieldFilterField) {
	mfs := make([]*FieldFilterField, 0)
	mfs = append(mfs, newFieldFilterField(nil, "NoCount", cde.Type(false)))
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldFilterFields(w.entity.Field(i))...)
	}
	strT := cde.Struct(fmt.Sprintf("%sFilter%s%s", w.entityName, w.OutFilterSuffix, w.OutTypeSuffix))
	strT.AddFields(fieldFilterFieldsToGocoder(newFieldFilterField(nil, "Ands", strT.GetType().TackPtr().Slice())))
	strT.AddFields(fieldFilterFieldsToGocoder(newFieldFilterField(nil, "Ors", strT.GetType().TackPtr().Slice())))
	strT.AddFields(fieldFilterFieldsToGocoder(mfs...))
	return strT, mfs
}

func (w *RepositoryWriter) GetEntityStruct() gocoder.Struct {
	mfs := make([]gocoder.Field, 0)
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, w.entity.Field(i))
	}
	strT := cde.Struct(fmt.Sprintf("%s"))
	strT.AddFields(mfs)
	return strT
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
	mfs = append(mfs, newFieldUpdaterField(nil, "JustDelete", cde.Type(true)))
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldUpdaterFields(w.entity.Field(i))...)
	}

	strT := cde.Struct(fmt.Sprintf("%sUpdater%s%s", w.entityName, w.OutUpdaterSuffix, w.OutTypeSuffix), fieldUpdaterFieldsToGocoder(mfs)...)
	return strT, mfs
}

func (w *RepositoryWriter) GetUpdaterTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	strT, mfs := w.GetUpdaterTypeStructCode()
	c.C(strT)
	c.C(getFieldUpdaterMethodToBSON(strT, mfs))

	return c
}

func (w *RepositoryWriter) GetSorterTypeStructCode() (gocoder.Struct, []*FieldSorterField) {
	mfs := make([]*FieldSorterField, 0)
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldSorterFields(w.entity.Field(i))...)
	}

	strT := cde.Struct(fmt.Sprintf("%sSorter%s%s", w.entityName, w.OutSorterSuffix, w.OutTypeSuffix), fieldSorterFieldsToGocoder(mfs)...)
	return strT, mfs
}

func (w *RepositoryWriter) GetSelectorTypeStructCode() (gocoder.Struct, []*FieldSelectorField) {
	mfs := make([]*FieldSelectorField, 0)
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldSelectorFields(w.entity.Field(i))...)
	}

	strT := cde.Struct(fmt.Sprintf("%sSelector%s%s", w.entityName, w.OutSelectorSuffix, w.OutTypeSuffix), fieldSelectorFieldsToGocoder(mfs)...)
	return strT, mfs
}

func (w *RepositoryWriter) GetSorterTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	strT, mfs := w.GetSorterTypeStructCode()
	c.C(strT)
	c.C(getFieldSorterMethodToBSON(strT, mfs))

	return c
}

func (w *RepositoryWriter) GetSelectorTypeCode() gocoder.Code {
	c := gocoder.NewCode()
	strT, mfs := w.GetSelectorTypeStructCode()
	c.C(strT)
	c.C(getFieldSelectorMethodToBSON(strT, mfs))

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

func (w *RepositoryWriter) NewEntityTmplContext() *EntityTmplContext {
	return &EntityTmplContext{
		w: w,
	}
}

func (w *RepositoryWriter) GetEntityRepositoryCodeFromTmpl(tmplPath string, ctx *EntityTmplContext) (gocoder.Code, error) {
	c := gocoder.NewCode()
	code, err := gocoder.TemplateFromFile(tmplPath, ctx, nil)
	if err != nil {
		return nil, err
	}
	c.C(code)
	return c, nil
}

func (w *RepositoryWriter) GetEntityRepositoryStructCode() gocoder.Code {
	return w.GetEntityRepositoryCode(w.GetFilterTypeStructCodeStruct(), w.GetUpdaterTypeStructCodeStruct())
}

func (w *RepositoryWriter) GetEntityRepositoryInterfaceCode() gocoder.Code {
	c := gocoder.NewCode()

	mfs := make([]gocoder.Func, 0)
	mfs = append(mfs, cde.Func("InitDB", []gocoder.Arg{cde.Arg("ctx", cdt.Context())}, nil))
	mfs = append(mfs, cde.Func("Insert", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.ArgVar("docs", cde.TypeD(w.entityPkg, w.entityName))}, []gocoder.Type{cde.TypeError()}))
	mfs = append(mfs, cde.Func("Update", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("filter", cde.TypeD("", fmt.Sprintf("%sFilter", w.entityName))), cde.Arg("updater", cde.TypeD("", fmt.Sprintf("%sUpdater", w.entityName)))}, []gocoder.Type{cde.TypeError()}))
	mfs = append(mfs, cde.Func("Get", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("id", cdt.String())}, []gocoder.Type{cde.TypeD(w.entityPkg, w.entityName).TackPtr(), cde.TypeError()}))
	mfs = append(mfs, cde.Func("GetBatch", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("ids", cdt.StringSlice())}, []gocoder.Type{cde.TypeD(w.entityPkg, w.entityName).TackPtr().Slice(), cde.TypeError()}))
	mfs = append(mfs, cde.Func("Query", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("filter", cde.TypeD("", fmt.Sprintf("%sFilter", w.entityName))), cde.Arg("skipLen", cdt.Int()), cde.Arg("limitLen", cdt.Int())}, []gocoder.Type{cde.TypeD(w.entityPkg, w.entityName).TackPtr().Slice(), cdt.Int(), cde.TypeError()}))

	strT := cde.Interface(fmt.Sprintf("Base%sRepository", w.entityName), mfs...)
	c.C(strT)
	return c
}
