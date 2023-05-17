package context

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/liasece/gocoder"
)

type TypeTmplContext struct {
	*TmplContext
	gocoder.Type
	fields          []*EntityStructFieldTmplContext
	fieldFieldsTmpl *FieldsTmplContext
}

func NewTypeTmplContextList(ctx *TmplContext, methods []gocoder.Type) []*TypeTmplContext {
	res := make([]*TypeTmplContext, 0, len(methods))
	for _, m := range methods {
		res = append(res, NewTypeTmplContext(ctx, m))
	}
	return res
}

func NewTypeTmplContext(ctx *TmplContext, typ gocoder.Type) *TypeTmplContext {
	var fs []gocoder.Field
	if typ != nil && typ.Kind() == reflect.Struct {
		fs = typ.GetFields()
	}
	fields := NewEntityStructFieldTmplContextList(ctx, fs)
	res := &TypeTmplContext{
		TmplContext:     ctx,
		Type:            typ,
		fields:          fields,
		fieldFieldsTmpl: nil,
	}
	{
		fieldFields := make([]IField, 0, len(fields))
		for _, v := range fields {
			fieldFields = append(fieldFields, v)
		}
		fieldFieldsTmpl := NewFieldsTmplContext(ctx, fieldFields)
		res.fieldFieldsTmpl = fieldFieldsTmpl
	}
	return res
}

func (e *TypeTmplContext) Elem() *TypeTmplContext {
	if next := e.Type.GetNext(); next != nil {
		return NewTypeTmplContext(e.TmplContext, next)
	}
	return NewTypeTmplContext(e.TmplContext, e.Type.Elem())
}

func (e *TypeTmplContext) FinalElem() *TypeTmplContext {
	elem := e
	for elem.KindIsPointer() || elem.KindIsSlice() || elem.KindIsArray() || elem.KindIsChan() || elem.KindIsMap() {
		elem = elem.Elem()
	}
	return elem
}

func unEnumType(typ gocoder.Type) gocoder.Type {
	switch typ.Kind() {
	case reflect.String:
		if typ.String() != "string" {
			return gocoder.NewTypeI("string")
		}
	case reflect.Int:
		if typ.String() != "int" {
			return gocoder.NewTypeI(int(0))
		}
	case reflect.Int8:
		if typ.String() != "int8" {
			return gocoder.NewTypeI(int8(0))
		}
	case reflect.Int16:
		if typ.String() != "int16" {
			return gocoder.NewTypeI(int16(0))
		}
	case reflect.Int32:
		if typ.String() != "int32" {
			return gocoder.NewTypeI(int32(0))
		}
	case reflect.Int64:
		if typ.String() != "int64" {
			return gocoder.NewTypeI(int64(0))
		}
	case reflect.Uint:
		if typ.String() != "uint" {
			return gocoder.NewTypeI(uint(0))
		}
	case reflect.Uint8:
		if typ.String() != "uint8" {
			return gocoder.NewTypeI(uint8(0))
		}
	case reflect.Uint16:
		if typ.String() != "uint16" {
			return gocoder.NewTypeI(uint16(0))
		}
	case reflect.Uint32:
		if typ.String() != "uint32" {
			return gocoder.NewTypeI(uint32(0))
		}
	case reflect.Uint64:
		if typ.String() != "uint64" {
			return gocoder.NewTypeI(uint64(0))
		}
	case reflect.Ptr:
		newType := unEnumType(typ.Elem())
		return newType.TackPtr()
	case reflect.Slice:
		newType := unEnumType(typ.Elem())
		return newType.Slice()
	default:
		return typ
	}
	return typ
}

func (e *TypeTmplContext) UnEnum() *TypeTmplContext {
	return NewTypeTmplContext(e.TmplContext, unEnumType(e.Type))
}

func (e *TypeTmplContext) UnEnumString() string {
	return e.UnEnum().String()
}

func (e *TypeTmplContext) ExternalTypeString() string {
	str := e.Type.Name()
	finalElem := e.FinalElem()
	finalElemPkg := finalElem.PackageInReference()
	if finalElemPkg != "" {
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

// docReg like `@description\s+(.*)` group like 1, doc like `@description xxx`, return `xxx`
func (e *TypeTmplContext) GetDocByReg(docReg string, group int) string {
	reg := regexp.MustCompile(docReg)
	for _, note := range e.Type.Notes() {
		if ss := reg.FindStringSubmatch(note.GetContent()); len(ss) > 0 {
			return ss[group]
		}
	}
	return ""
}

func (e *TypeTmplContext) Doc() string {
	resList := make([]string, 0)
	for _, note := range e.Type.Notes() {
		resList = append(resList, note.GetContent())
	}
	return strings.Join(resList, "\n")
}

func (e *TypeTmplContext) DocLinesTrimAndJoin(joinStr string) string {
	return docLinesTrimAndJoin(e.Doc(), joinStr)
}

func (e *TypeTmplContext) FieldsGraphqlDefinition() string {
	return e.fieldFieldsTmpl.GraphqlDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Type().Name() != "Context"
	})
}

func (e *TypeTmplContext) FieldsProtoBuffDefinition() string {
	return e.fieldFieldsTmpl.ProtoBuffDefinitionFilterFunc(func(i IField) bool {
		return i.Type().Name() != "error" && i.Type().Name() != "Context"
	})
}
