package cmd

import (
	"github.com/liasece/go-mate/config"
	ccontext "github.com/liasece/go-mate/context"
	"github.com/liasece/go-mate/utils"
	"github.com/liasece/gocoder"
	"github.com/liasece/log"
)

func generateMethods(entityCfg *config.Entity) {
	if entityCfg.CodeType == nil {
		return
	}
	cs := entityCfg.CodeType.(gocoder.Code).GetCodes()
	methods := make([]gocoder.Func, 0)
	for _, c := range cs {
		methods = append(methods, c.(gocoder.Func))
	}
	for _, tmpl := range entityCfg.Tmpl {
		generateMethodsTmplItem(entityCfg, tmpl, methods)
	}
}

func generateMethodsTmplItem(entityCfg *config.Entity, tmpl *config.TmplItem, methods []gocoder.Func) {
	toFile, err := utils.TemplateRaw(tmpl.To, ccontext.NewMethodsTmplContext(ccontext.NewTmplContext(tmpl, entityCfg), methods))
	if err != nil {
		log.Fatal("generateMethods TemplateRaw error", log.ErrorField(err))
		return
	}
	tmplCtx := ccontext.NewMethodsTmplContext(ccontext.NewTmplContext(tmpl, entityCfg), methods)
	generateEntityTmplToFile(tmplCtx, entityCfg.CodeName, toFile, tmpl)
}
