package render

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type MultipleHTML struct {
	Templates map[string]*template.Template
}

// New instance
func NewMultipleHTML() MultipleHTML {
	return MultipleHTML{
		Templates: make(map[string]*template.Template),
	}
}

// Add new template
func (r MultipleHTML) Add(name string, tmpl *template.Template) {
	if tmpl == nil {
		panic("template can not be nil")
	}
	if len(name) == 0 {
		panic("template name cannot be empty")
	}
	if _, ok := r.Templates[name]; ok {
		panic(fmt.Sprintf("template %s already exists", name))
	}
	r.Templates[name] = tmpl
}

// AddFromFilesFuncs supply add template from file callback func
func (r MultipleHTML) AddFromFilesFuncs(name string, funcMap template.FuncMap, files ...string) *template.Template {
	tname := filepath.Base(files[0])
	tmpl := template.Must(template.New(tname).Funcs(funcMap).ParseFiles(files...))
	r.Add(name, tmpl)
	return tmpl
}

// Instance supply render string
func (r MultipleHTML) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Templates[name],
		Data:     data,
	}
}
