package template

import (
	"github.com/vuongabc92/octocv/render"
	"html/template"
)

type Engine struct {
	delims     render.Delims
	HTMLRender render.HTMLRender
	FuncMap    template.FuncMap
}

func NewTemplateEngine() *Engine {
	return &Engine{
		delims:  render.Delims{Left: "{{", Right: "}}"},
		FuncMap: template.FuncMap{},
	}
}

// SetHTMLTemplate associate a template with HTML renderer.
func (engine *Engine) SetHTMLTemplate(t *template.Template) {
	engine.HTMLRender = render.HTMLProduction{Template: t.Funcs(engine.FuncMap)}
}
