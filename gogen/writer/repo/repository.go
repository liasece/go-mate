package repo

import (
	"fmt"

	"github.com/liasece/go-mate/config"
	"github.com/liasece/gocoder"
	"github.com/liasece/gocoder/cde"
	"github.com/liasece/gocoder/cdt"
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

		entityPkg:         "",
		serviceName:       "",
		Filter:            nil,
		Updater:           nil,
		Sorter:            nil,
		Selector:          nil,
		OutTypeSuffix:     "",
		OutFilterSuffix:   "",
		OutUpdaterSuffix:  "",
		OutSorterSuffix:   "",
		OutSelectorSuffix: "",
		EntityCfg:         nil,
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

		Filter:    nil,
		Updater:   nil,
		Sorter:    nil,
		Selector:  nil,
		EntityCfg: nil,
	}
}

func (w *RepositoryWriter) GetFilterTypeStructCodeStruct() gocoder.Struct {
	if w.Filter != nil {
		return w.Filter
	}
	res, _ := w.GetFilterTypeStructCode()
	return res
}

func (w *RepositoryWriter) Entity() gocoder.Type {
	return w.entity
}

func (w *RepositoryWriter) EntityName() string {
	return w.entityName
}

func (w *RepositoryWriter) ServiceName() string {
	return w.serviceName
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
	strT := cde.Struct(w.entityName)
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

type TypeCfg struct {
	addSlicePBEmpty bool
}

type TypeOpt func(*TypeCfg)

func TypeOptAddSlicePBEmpty(v bool) TypeOpt {
	return func(cfg *TypeCfg) {
		cfg.addSlicePBEmpty = v
	}
}

func (w *RepositoryWriter) GetUpdaterTypeStructCode(opts ...TypeOpt) (gocoder.Struct, []*FieldUpdaterField) {
	mfs := make([]*FieldUpdaterField, 0)
	mfs = append(mfs, newFieldUpdaterField(nil, "JustDelete", cde.Type(true)))
	for i := 0; i < w.entity.NumField(); i++ {
		mfs = append(mfs, getFieldUpdaterFields(w.entity.Field(i), opts...)...)
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

func (w *RepositoryWriter) GetEntityRepositoryInterfaceCode() gocoder.Code {
	c := gocoder.NewCode()

	mfs := make([]gocoder.Func, 0)
	mfs = append(mfs, cde.Func("InitDB", []gocoder.Arg{cde.Arg("ctx", cdt.Context())}, nil))
	mfs = append(mfs, cde.Func("Insert", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.ArgVar("docs", cde.TypeD(w.entityPkg, w.entityName))}, []gocoder.Arg{cde.Arg("", cde.TypeError())}))
	mfs = append(mfs, cde.Func("Update", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("filter", cde.TypeD("", fmt.Sprintf("%sFilter", w.entityName))), cde.Arg("updater", cde.TypeD("", fmt.Sprintf("%sUpdater", w.entityName)))}, []gocoder.Arg{cde.Arg("", cde.TypeError())}))
	mfs = append(mfs, cde.Func("Get", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("id", cdt.String())}, []gocoder.Arg{cde.Arg("", cde.TypeD(w.entityPkg, w.entityName).TackPtr()), cde.Arg("", cde.TypeError())}))
	mfs = append(mfs, cde.Func("GetBatch", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("ids", cdt.StringSlice())}, []gocoder.Arg{cde.Arg("", cde.TypeD(w.entityPkg, w.entityName).TackPtr().Slice()), cde.Arg("", cde.TypeError())}))
	mfs = append(mfs, cde.Func("Query", []gocoder.Arg{cde.Arg("ctx", cdt.Context()), cde.Arg("filter", cde.TypeD("", fmt.Sprintf("%sFilter", w.entityName))), cde.Arg("skipLen", cdt.Int()), cde.Arg("limitLen", cdt.Int())}, []gocoder.Arg{cde.Arg("", cde.TypeD(w.entityPkg, w.entityName).TackPtr().Slice()), cde.Arg("", cdt.Int()), cde.Arg("", cde.TypeError())}))

	strT := cde.Interface(fmt.Sprintf("Base%sRepository", w.entityName), mfs...)
	c.C(strT)
	return c
}
