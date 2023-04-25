package context

import (
	"sort"

	"github.com/liasece/gocoder"
)

type MethodsTmplSortByName []*MethodTmplContext

func (a MethodsTmplSortByName) Len() int           { return len(a) }
func (a MethodsTmplSortByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MethodsTmplSortByName) Less(i, j int) bool { return a[i].Name() < a[j].Name() }

type MethodsTmplContext struct {
	*TmplContext
	methods []*MethodTmplContext
}

func NewMethodsTmplContext(ctx *TmplContext, methods []gocoder.Func) *MethodsTmplContext {
	list := NewMethodTmplContextList(ctx, methods)
	sort.Sort(MethodsTmplSortByName(list))
	return &MethodsTmplContext{
		TmplContext: ctx,
		methods:     list,
	}
}

func (c *MethodsTmplContext) Methods() []*MethodTmplContext {
	return c.methods
}

func (c *MethodsTmplContext) HaveMethods() bool {
	return len(c.methods) > 0
}

func (c *MethodsTmplContext) FindMethods(nameReg string) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0)
	for _, m := range c.methods {
		if m.IsNameReg(nameReg) {
			res = append(res, m)
		}
	}
	return res
}

func (c *MethodsTmplContext) HaveMethodsByName(nameReg string) bool {
	return len(c.FindMethods(nameReg)) > 0
}

func (c *MethodsTmplContext) FindMethodsNot(nameReg string) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0)
	for _, m := range c.methods {
		if !m.IsNameReg(nameReg) {
			res = append(res, m)
		}
	}
	return res
}

func (c *MethodsTmplContext) HaveMethodsNot(docReg string) bool {
	return len(c.FindMethodsNot(docReg)) > 0
}

func (c *MethodsTmplContext) FindMethodsByDoc(docReg string) []*MethodTmplContext {
	res := make([]*MethodTmplContext, 0)
	for _, m := range c.methods {
		if m.IsDocReg(docReg) {
			res = append(res, m)
		}
	}
	return res
}

func (c *MethodsTmplContext) HaveMethodsByDoc(docReg string) bool {
	return len(c.FindMethodsByDoc(docReg)) > 0
}
