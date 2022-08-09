package repo

import (
	"reflect"
	"strings"

	"github.com/liasece/go-mate/src/gogen/utils"
)

type EntityTmplContext struct {
	utils.TmplUtilsFunc
	w *RepositoryWriter
}

func (e *EntityTmplContext) EntityName() string {
	return e.w.entityName
}

func (e *EntityTmplContext) ServiceName() string {
	return e.w.serviceName
}

func (e *EntityTmplContext) ServiceNameTitle() string {
	s := e.ServiceName()
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (e *EntityTmplContext) GetTagOn(filterReg string, targetTag string) string {
	fieldNum := e.findFieldNumByTagOn(filterReg)
	if fieldNum < 0 {
		return ""
	}
	if targetTag == "" {
		return e.w.entity.Field(fieldNum).GetName()
	}
	return reflect.StructTag(e.w.entity.Field(fieldNum).GetTag()).Get(targetTag)
}

func (e *EntityTmplContext) GetTypeByTagOn(filterReg string) string {
	fieldNum := e.findFieldNumByTagOn(filterReg)
	if fieldNum < 0 {
		return ""
	}
	return e.w.entity.Field(fieldNum).GetType().String()
}

func (e *EntityTmplContext) GetType(filedName string) string {
	for i := 0; i < e.w.entity.NumField(); i++ {
		if e.w.entity.Field(i).GetName() == filedName {
			return e.w.entity.Field(i).GetType().String()
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
	for i := 0; i < e.w.entity.NumField(); i++ {
		t := reflect.StructTag(e.w.entity.Field(i).GetTag())
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
