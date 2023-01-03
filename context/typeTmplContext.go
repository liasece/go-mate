package context

import (
	"reflect"
	"strings"

	"github.com/liasece/gocoder"
)

type TypeTmplContext struct {
	*TmplContext
	gocoder.Type
}

func NewTypeTmplContextList(ctx *TmplContext, methods []gocoder.Type) []*TypeTmplContext {
	res := make([]*TypeTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewTypeTmplContext(ctx, m))
	}
	return res
}

func NewTypeTmplContext(ctx *TmplContext, typ gocoder.Type) *TypeTmplContext {
	return &TypeTmplContext{
		TmplContext: ctx,
		Type:        typ,
	}
}

func (e *TypeTmplContext) Elem() *TypeTmplContext {
	if next := e.Type.GetNext(); next != nil {
		return &TypeTmplContext{
			TmplContext: e.TmplContext,
			Type:        next,
		}
	}
	return &TypeTmplContext{
		TmplContext: e.TmplContext,
		Type:        e.Type.Elem(),
	}
}

func (e *TypeTmplContext) FinalElem() *TypeTmplContext {
	elem := e
	for elem.KindIsPointer() || elem.KindIsSlice() || elem.KindIsArray() || elem.KindIsChan() || elem.KindIsMap() || elem.Type.GetNext() != nil {
		elem = elem.Elem()
	}
	return elem
}

func (e *TypeTmplContext) ExternalTypeString() string {
	str := e.Type.Name()
	finalElem := e.FinalElem()
	finalElemPkg := finalElem.PackageInReference()
	if finalElem.KindIsStruct() && finalElemPkg != "" {
		str = strings.ReplaceAll(str, finalElem.Name(), finalElemPkg+"."+finalElem.Name())
	}
	return str
}

func (e *TypeTmplContext) KindIsNumber() bool {
	return e.Type.Kind() >= reflect.Int && e.Type.Kind() <= reflect.Float64
}

func (e *TypeTmplContext) KindIsBool() bool {
	return e.Type.Kind() == reflect.Bool
}

func (e *TypeTmplContext) KindIsInt() bool {
	return e.Type.Kind() == reflect.Int
}

func (e *TypeTmplContext) KindIsInt8() bool {
	return e.Type.Kind() == reflect.Int8
}

func (e *TypeTmplContext) KindIsInt16() bool {
	return e.Type.Kind() == reflect.Int16
}

func (e *TypeTmplContext) KindIsInt32() bool {
	return e.Type.Kind() == reflect.Int32
}

func (e *TypeTmplContext) KindIsInt64() bool {
	return e.Type.Kind() == reflect.Int64
}

func (e *TypeTmplContext) KindIsUint() bool {
	return e.Type.Kind() == reflect.Uint
}

func (e *TypeTmplContext) KindIsUint8() bool {
	return e.Type.Kind() == reflect.Uint8
}

func (e *TypeTmplContext) KindIsUint16() bool {
	return e.Type.Kind() == reflect.Uint16
}

func (e *TypeTmplContext) KindIsUint32() bool {
	return e.Type.Kind() == reflect.Uint32
}

func (e *TypeTmplContext) KindIsUint64() bool {
	return e.Type.Kind() == reflect.Uint64
}

func (e *TypeTmplContext) KindIsUintptr() bool {
	return e.Type.Kind() == reflect.Uintptr
}

func (e *TypeTmplContext) KindIsFloat32() bool {
	return e.Type.Kind() == reflect.Float32
}

func (e *TypeTmplContext) KindIsFloat64() bool {
	return e.Type.Kind() == reflect.Float64
}

func (e *TypeTmplContext) KindIsComplex64() bool {
	return e.Type.Kind() == reflect.Complex64
}

func (e *TypeTmplContext) KindIsComplex128() bool {
	return e.Type.Kind() == reflect.Complex128
}

func (e *TypeTmplContext) KindIsArray() bool {
	return e.Type.Kind() == reflect.Array
}

func (e *TypeTmplContext) KindIsChan() bool {
	return e.Type.Kind() == reflect.Chan
}

func (e *TypeTmplContext) KindIsFunc() bool {
	return e.Type.Kind() == reflect.Func
}

func (e *TypeTmplContext) KindIsInterface() bool {
	return e.Type.Kind() == reflect.Interface
}

func (e *TypeTmplContext) KindIsMap() bool {
	return e.Type.Kind() == reflect.Map
}

func (e *TypeTmplContext) KindIsPointer() bool {
	return e.Type.Kind() == reflect.Pointer
}

func (e *TypeTmplContext) KindIsSlice() bool {
	return e.Type.Kind() == reflect.Slice
}

func (e *TypeTmplContext) KindIsString() bool {
	return e.Type.Kind() == reflect.String
}

func (e *TypeTmplContext) KindIsStruct() bool {
	return e.Type.Kind() == reflect.Struct
}

func (e *TypeTmplContext) KindIsUnsafePointer() bool {
	return e.Type.Kind() == reflect.UnsafePointer
}