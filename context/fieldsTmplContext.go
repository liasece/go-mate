package context

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/liasece/go-mate/utils"
)

type IField interface {
	Type() *TypeTmplContext
	Name() string

	GraphqlDefinition() string
	Doc() string
	GetTmplContext() *TmplContext
	ProtoBuffDefinition(fieldIndex int) string
}

type FieldsTmplContext struct {
	*TmplContext
	docReader func(fieldName string) (doc string)
	fields    []IField
}

func NewFieldsTmplContext(ctx *TmplContext, fields []IField) *FieldsTmplContext {
	return &FieldsTmplContext{
		TmplContext: ctx,
		docReader:   nil,
		fields:      fields,
	}
}

func (c *FieldsTmplContext) GetFieldDoc(fieldName string) string {
	var docs []string
	for _, field := range c.fields {
		if field.Name() == fieldName {
			docs = append(docs, c.getFieldDoc(field))
		}
	}
	return strings.Join(docs, "\n")
}

func (c *FieldsTmplContext) getFieldDoc(field IField) string {
	var docs []string

	if doc := field.Doc(); doc != "" {
		docs = append(docs, doc)
	}
	if c.docReader != nil {
		if doc := c.docReader(field.Name()); doc != "" {
			docs = append(docs, doc)
		}
	}

	return strings.Join(docs, "\n")
}

func (c *FieldsTmplContext) GraphqlDefinition() string {
	return c.GraphqlDefinitionFilterFunc(nil)
}

func docLinesTrimAndJoin(doc string, joinStr string) string {
	ss := strings.Split(doc, "\n")
	for i, s := range ss {
		ss[i] = strings.TrimSpace(s)
	}
	return strings.Join(ss, joinStr)
}

func (c *FieldsTmplContext) GraphqlDefinitionFilterFunc(filter func(IField) bool) string {
	res := ""
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		argType := NewTypeTmplContext(c.TmplContext, arg.Type().UnPtr())
		if argType.IsStruct() && strings.HasSuffix(argType.Name(), "Input") {
			res += argType.FieldsGraphqlDefinition()
			continue
		}
		definitionStr := arg.GraphqlDefinition()
		if definitionStr == "" {
			continue
		}
		{
			// add doc
			if doc := c.getFieldDoc(arg); doc != "" {
				res += fmt.Sprintf("  \"\"\"\n  %s\n  \"\"\"\n", docLinesTrimAndJoin(doc, "\n  "))
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) ProtoBuffDefinition() string {
	return c.ProtoBuffDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) debugString() string {
	res := ""
	for _, v := range c.fields {
		if v == nil {
			res += "nil\n"
		} else {
			res += fmt.Sprintf("%s\n", v.Name())
		}
	}
	return res
}

func (c *FieldsTmplContext) ProtoBuffDefinitionFilterFunc(filter func(IField) bool) string {
	res := ""
	argIndex := 1
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		if arg == nil || arg.Type() == nil {
			panic(c.debugString())
		}
		argType := NewTypeTmplContext(arg.GetTmplContext(), arg.Type().UnPtr())
		if argType.IsStruct() && strings.HasSuffix(argType.Name(), "Input") {
			res += argType.FieldsProtoBuffDefinition()
			continue
		}
		definitionStr := arg.ProtoBuffDefinition(argIndex)
		if definitionStr == "" {
			continue
		}
		{
			// add doc
			if doc := c.getFieldDoc(arg); doc != "" {
				res += fmt.Sprintf("  //%s\n", doc)
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
		argIndex++
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) GRPCCallGoDefinition(reqValueName string) string {
	return c.GRPCCallGoDefinitionFilterFunc(nil, reqValueName)
}

func (c *FieldsTmplContext) GRPCCallGoDefinitionFilterFunc(filter func(IField) bool, reqValueName string) string {
	res := []string{}
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		if arg.Type().Name() == "Context" {
			res = append(res, "ctx")
			continue
		}
		res = append(res, reqValueName+"."+utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GoDefinition() string {
	return c.GoDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) GoDefinitionFilterFunc(filter func(IField) bool) string {
	res := []string{}
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		res = append(res, fmt.Sprintf("%s %s", arg.Name(), arg.Type().Type.String()))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GraphqlGoDefinition() string {
	return c.GraphqlGoDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) GraphqlGoDefinitionFilterFunc(filter func(IField) bool) string {
	res := []string{}
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		typ := arg.Type().Type.String()
		if pkg := arg.Type().Type.PackageInReference(); pkg != "" {
			typ = pkg + "." + typ
		}
		typ = strings.ReplaceAll(typ, "int32", "int")
		typ = strings.ReplaceAll(typ, "int64", "int")
		res = append(res, fmt.Sprintf("%s %s", arg.Name(), typ))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) CallGRPCDefinition() string {
	return c.CallGRPCDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) CallGRPCDefinitionFilterFunc(filter func(IField) bool) string {
	res := []string{}
	for _, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		value := arg.Name()
		if arg.Type().KindIsNumber() {
			value = "int64(" + value + ")"
		}
		if arg.Type().Kind() == reflect.Slice && arg.Type().Elem().KindIsNumber() {
			value = "func() []int64 {res := []int64{}; for _,v := range " + value + " { res = append(res, int64(v)) }; return res } ()"
		}
		res = append(res, fmt.Sprintf("%s: %s", utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())), value))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) CallGoDefinition() string {
	return c.CallGoDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) CallGoDefinitionFilterFunc(filter func(IField) bool) string {
	res := []string{}
	for i, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		goName := arg.Name()
		if goName == "" {
			if arg.Type().Name() == "error" {
				goName = "err"
			} else {
				goName = fmt.Sprintf("v%d", i)
			}
		}
		res = append(res, goName)
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GRPCCallGoResponseDefinition() string {
	return c.GRPCCallGoResponseDefinitionFilterFunc(nil)
}

func (c *FieldsTmplContext) GRPCCallGoResponseDefinitionFilterFunc(filter func(IField) bool) string {
	res := []string{}
	for i, arg := range c.fields {
		if filter != nil && !filter(arg) {
			continue
		}
		goName := arg.Name()
		if goName == "" {
			if arg.Type().Name() == "error" {
				goName = "err"
			} else {
				goName = fmt.Sprintf("v%d", i)
			}
		}
		res = append(res, fmt.Sprintf("%s: %s", utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())), goName))
	}
	return strings.Join(res, ", \n")
}
