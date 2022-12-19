package context

import (
	"context"
	"reflect"
	"strings"

	"github.com/liasece/go-mate/src/gogen/writer/repo"
	"github.com/liasece/go-mate/src/utils"
)

type EntityTmplContext struct {
	*TmplContext
	w *repo.RepositoryWriter
}

func NewEntityTmplContext(ctx *TmplContext, w *repo.RepositoryWriter) *EntityTmplContext {
	return &EntityTmplContext{
		TmplContext: ctx,
		w:           w,
	}
}

func (e *EntityTmplContext) EntityName() string {
	return e.w.EntityName()
}

func (e *EntityTmplContext) ServiceName() string {
	return e.w.ServiceName()
}

func (e *EntityTmplContext) ServiceNameTitle() string {
	s := e.ServiceName()
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (e *EntityTmplContext) EntityStruct() *EntityStructTmplContext {
	return &EntityStructTmplContext{
		w:      e.w,
		Struct: e.w.GetEntityStruct(),
	}
}

func (e *EntityTmplContext) Sorter() *EntityStructTmplContext {
	return &EntityStructTmplContext{
		w:      e.w,
		Struct: e.w.Sorter,
	}
}

func (e *EntityTmplContext) Selector() *EntityStructTmplContext {
	return &EntityStructTmplContext{
		w:      e.w,
		Struct: e.w.Selector,
	}
}

func (e *EntityTmplContext) Filter() *EntityStructTmplContext {
	return &EntityStructTmplContext{
		w:      e.w,
		Struct: e.w.Filter,
	}
}

func (e *EntityTmplContext) Updater() *EntityStructTmplContext {
	return &EntityStructTmplContext{
		w:      e.w,
		Struct: e.w.Updater,
	}
}

func (e *EntityTmplContext) Env(k1 string, k2 string) string {
	if e.w.EntityCfg != nil {
		if _, ok := e.w.EntityCfg.Env[k1]; ok {
			if v, ok := e.w.EntityCfg.Env[k1][k2]; ok {
				return v
			}
		}
	}
	return ""
}

func (e *EntityTmplContext) EnvOr(k1 string, k2 string, or string) string {
	if v := e.Env(k1, k2); v != "" {
		return v
	}
	return or
}

func (e *EntityTmplContext) GetTagOn(filterReg string, targetTag string) string {
	fieldNum := e.findFieldNumByTagOn(filterReg)
	if fieldNum < 0 {
		return ""
	}
	if targetTag == "" {
		return e.w.Entity().Field(fieldNum).GetName()
	}
	return reflect.StructTag(e.w.Entity().Field(fieldNum).GetTag()).Get(targetTag)
}

func (e *EntityTmplContext) GetTypeByTagOn(filterReg string) string {
	fieldNum := e.findFieldNumByTagOn(filterReg)
	if fieldNum < 0 {
		return ""
	}
	return e.w.Entity().Field(fieldNum).GetType().String()
}

func (e *EntityTmplContext) GetType(filedName string) string {
	for i := 0; i < e.w.Entity().NumField(); i++ {
		if e.w.Entity().Field(i).GetName() == filedName {
			return e.w.Entity().Field(i).GetType().String()
		}
	}
	return ""
}

func (e *EntityTmplContext) EntityGrpcSubPkg() string {
	if e.w.EntityCfg == nil {
		return ""
	}
	return e.w.EntityCfg.GrpcSubPkg
}

func (e *EntityTmplContext) findFieldNumByTagOn(filterReg string) int {
	filterSS := strings.Split(filterReg, ":")
	filterTag := filterSS[0]
	filterValue := ""
	if len(filterSS) > 1 {
		filterValue = filterSS[1]
	}
	for i := 0; i < e.w.Entity().NumField(); i++ {
		t := reflect.StructTag(e.w.Entity().Field(i).GetTag())
		find := false
		if value := t.Get(filterTag); value != "" {
			if filterValue == "" {
				find = true
			} else {
				values := strings.Split(value, ",")
				for _, v := range values {
					if v == filterValue {
						find = true
						break
					}
				}
			}
		}
		if !find {
			continue
		}
		return i
	}
	return -1
}

// get gorm indexes
func (e *EntityTmplContext) GormIndexes() []*utils.Index {
	return utils.GormIndexes(context.Background(), e.w.Entity())
}
