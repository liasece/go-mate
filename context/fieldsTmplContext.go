package context

import (
	"fmt"
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
			if doc := field.Doc(); doc != "" {
				docs = append(docs, doc)
			}
			if c.docReader != nil {
				if doc := c.docReader(fieldName); doc != "" {
					docs = append(docs, doc)
				}
			}
		}
	}
	return strings.Join(docs, "\n")
}

func (c *FieldsTmplContext) GraphqlArgsDefinition() string {
	res := ""
	for _, arg := range c.fields {
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
			if doc := c.GetFieldDoc(arg.Name()); doc != "" {
				res += fmt.Sprintf("  \"\"\"\n%s\n\"\"\"\n", doc)
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) GraphqlReturnsDefinition() string {
	res := ""
	for _, arg := range c.fields {
		definitionStr := arg.GraphqlDefinition()
		if definitionStr == "" {
			continue
		}
		{
			// add doc
			if doc := c.GetFieldDoc(arg.Name()); doc != "" {
				res += fmt.Sprintf("  \"\"\"\n%s\n\"\"\"\n", doc)
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) ProtoBuffArgsDefinition() string {
	res := ""
	argIndex := 1
	for _, arg := range c.fields {
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
			if doc := c.GetFieldDoc(arg.Name()); doc != "" {
				res += fmt.Sprintf("  //%s\n", doc)
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
		argIndex++
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) ProtoBuffReturnsDefinition() string {
	res := ""
	argIndex := 1
	for _, arg := range c.fields {
		definitionStr := arg.ProtoBuffDefinition(argIndex)
		if definitionStr == "" {
			continue
		}
		{
			// add doc
			if doc := c.GetFieldDoc(arg.Name()); doc != "" {
				res += fmt.Sprintf("  //%s\n", doc)
			}
		}
		res += fmt.Sprintf("  %s\n", definitionStr)
		argIndex++
	}
	return strings.TrimSpace(res)
}

func (c *FieldsTmplContext) GRPCCallGoArgsDefinition(reqValueName string) string {
	res := []string{}
	for _, arg := range c.fields {
		if arg.Type().Name() == "Context" {
			res = append(res, "ctx")
			continue
		}
		res = append(res, reqValueName+"."+utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GoArgsDefinition() string {
	res := []string{}
	for _, arg := range c.fields {
		res = append(res, fmt.Sprintf("%s %s", arg.Name(), arg.Type().Type.String()))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GraphqlGoArgsDefinition() string {
	res := []string{}
	for _, arg := range c.fields {
		typ := arg.Type().Type.String()
		if pkg := arg.Type().Type.PackageInReference(); pkg != "" {
			typ = pkg + "." + typ
		}
		res = append(res, fmt.Sprintf("%s %s", arg.Name(), typ))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) CallGRPCArgsDefinition() string {
	res := []string{}
	for _, arg := range c.fields {
		res = append(res, fmt.Sprintf("%s: %s", utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())), arg.Name()))
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) CallGoReturnsDefinition() string {
	res := []string{}
	for i, arg := range c.fields {
		name := arg.Name()
		if name == "" {
			if arg.Type().Name() == "error" {
				name = "err"
			} else {
				name = fmt.Sprintf("ret%d", i)
			}
		}
		res = append(res, name)
	}
	return strings.Join(res, ", ")
}

func (c *FieldsTmplContext) GRPCCallGoReturnsResponseDefinition() string {
	res := []string{}
	for i, arg := range c.fields {
		goName := arg.Name()
		if goName == "" {
			if arg.Type().Name() == "error" {
				goName = "err"
			} else {
				goName = fmt.Sprintf("ret%d", i)
			}
		}
		res = append(res, fmt.Sprintf("%s: %s", utils.SnakeStringToBigHump(utils.SnakeString(arg.Name())), goName))
	}
	return strings.Join(res, ", \n")
}
