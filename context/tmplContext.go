package context

import (
	"path/filepath"
	"strings"

	"github.com/liasece/go-mate/config"
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
)

func GetCodeFromTmpl(ctx interface{}, tmplPath string) (gocoder.Code, error) {
	code, err := utils.TemplateFromFile(tmplPath, ctx)
	if err != nil {
		return nil, err
	}
	c := gocoder.NewCode()
	c.C(code)
	return c, nil
}

type ITmplContext interface {
	Terminate() bool
	GetTerminate() bool
	FromFilePath(path string) string
	ToFilePath(path string) string
	GetToFilePath() string
	GetFromFilePath() string
}

type BaseTmplContext struct {
	terminate    bool
	toFilePath   string
	fromFilePath string
}

func (e *BaseTmplContext) Terminate() bool {
	e.terminate = true
	return e.terminate
}

func (e *BaseTmplContext) GetTerminate() bool {
	return e.terminate
}

func (e *BaseTmplContext) ToFilePath(path string) string {
	e.toFilePath = path
	return e.toFilePath
}

func (e *BaseTmplContext) FromFilePath(path string) string {
	e.fromFilePath = path
	return e.fromFilePath
}

func (e *BaseTmplContext) GetToFilePath() string {
	return e.toFilePath
}

func (e *BaseTmplContext) GetFromFilePath() string {
	return e.fromFilePath
}

type TmplContext struct {
	BaseTmplContext
	tmpl      *config.TmplItem
	EntityCfg *config.Entity
}

func NewTmplContext(tmpl *config.TmplItem, entityCfg *config.Entity) *TmplContext {
	return &TmplContext{
		EntityCfg:       entityCfg,
		BaseTmplContext: BaseTmplContext{terminate: false},
		tmpl:            tmpl,
	}
}

func (e *TmplContext) EntityName() string {
	return e.EntityCfg.Name
}

func (e *TmplContext) EntityCodeName() string {
	return e.EntityCfg.CodeName
}

func (e *TmplContext) ServiceName() string {
	return e.EntityCfg.Service
}

func (e *TmplContext) ConfigFileDir() string {
	return filepath.Dir(e.EntityCfg.ConfigFilePath)
}

func (e *TmplContext) ServiceNameTitle() string {
	s := e.ServiceName()
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func (e *TmplContext) EntityGrpcSubPkg() string {
	if e.EntityCfg == nil {
		return ""
	}
	return e.EntityCfg.GrpcSubPkg
}

func EntityEnv(entityCfg *config.Entity, k1 string, k2 string) string {
	if _, ok := entityCfg.Env[k1]; ok {
		if v, ok := entityCfg.Env[k1][k2]; ok {
			return v
		}
	}
	return ""
}

func (e *TmplContext) Env(k1 string, k2 string) string {
	if e.EntityCfg != nil {
		return EntityEnv(e.EntityCfg, k1, k2)
	}
	return ""
}

func (e *TmplContext) EnvOr(k1 string, k2 string, or string) string {
	if v := e.Env(k1, k2); v != "" {
		return v
	}
	return or
}

func (e *TmplContext) TmplFrom() string {
	return e.tmpl.From
}

func (e *TmplContext) TmplTo() string {
	return e.tmpl.To
}

func (e *TmplContext) TmplType() string {
	return string(e.tmpl.Type)
}

func (e *TmplContext) TmplMerge() bool {
	return e.tmpl.Merge
}

func (e *TmplContext) TmplOnlyCreate() bool {
	return e.tmpl.OnlyCreate
}

func (e *TmplContext) TmplOptEq(key string, value string) bool {
	return e.TmplOpt(key) == value
}

func (e *TmplContext) TmplOpt(key string) string {
	if res, ok := e.tmpl.Opt[key]; ok {
		return res
	}
	return ""
}
