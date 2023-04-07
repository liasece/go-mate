package context

import (
	"reflect"
	"strings"

	"github.com/liasece/go-mate/gogen/writer/repo"
	"github.com/liasece/go-mate/utils"
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
	fieldNumList := e.findFieldNumByTagOn(filterReg)
	if len(fieldNumList) == 0 {
		return ""
	}
	fieldNum := fieldNumList[0]
	if targetTag == "" {
		return e.w.Entity().Field(fieldNum).GetName()
	}
	return reflect.StructTag(e.w.Entity().Field(fieldNum).GetTag()).Get(targetTag)
}

func (e *EntityTmplContext) ListFieldByTag(filterReg string) []*EntityStructFieldTmplContext {
	fieldNumList := e.findFieldNumByTagOn(filterReg)
	res := []*EntityStructFieldTmplContext{}
	for _, fieldNum := range fieldNumList {
		res = append(res, &EntityStructFieldTmplContext{
			w:     e.w,
			Field: e.w.Entity().Field(fieldNum),
		})
	}
	return res
}

func (e *EntityTmplContext) GetTypeByTagOn(filterReg string) string {
	fieldNumList := e.findFieldNumByTagOn(filterReg)
	if len(fieldNumList) == 0 {
		return ""
	}
	fieldNum := fieldNumList[0]
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

func (e *EntityTmplContext) findFieldNumByTagOn(filterReg string) []int {
	filterSS := strings.Split(filterReg, ":")
	filterTag := filterSS[0]
	filterValue := ""
	if len(filterSS) > 1 {
		filterValue = filterSS[1]
	}
	res := []int{}
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
		if find {
			res = append(res, i)
		}
	}
	return res
}

// get gorm indexes
func (e *EntityTmplContext) GormIndexes() []*utils.Index {
	return utils.GormIndexes(e.w.Entity())
}
