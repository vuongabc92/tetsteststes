package render

import (
	"bytes"
	"html/template"
	"net/http"
)

// Delims represents a set of Left and Right delimiters for HTML template rendering.
type Delims struct {
	// Left delimiter, defaults to {{.
	Left string
	// Right delimiter, defaults to }}.
	Right string
}

type HTML struct {
	Template *template.Template
	Name     string
	Data     interface{}
}

// HTMLRender interface is to be implemented by HTMLProduction and HTMLDebug.
type HTMLRender interface {
	// Instance returns an HTML instance.
	Instance(string, interface{}) Render
}

// HTMLProduction contains template reference.
type HTMLProduction struct {
	Template *template.Template
}

var htmlContentType = []string{"text/html; charset=utf-8"}

// Instance (HTMLProduction) returns an HTML instance which it realizes Render interface.
func (r HTMLProduction) Instance(name string, data interface{}) Render {
	return HTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

// Render (HTML) executes template and writes its result with custom ContentType for response.
func (r HTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}

	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// Render (HTML) executes template and returns it as string.
func (r HTML) String() (str string, err error) {
	buf := new(bytes.Buffer)

	if r.Name == "" {
		err = r.Template.Execute(buf, r.Data)
	} else {
		err = r.Template.ExecuteTemplate(buf, r.Name, r.Data)
	}

	str = buf.String()
	return
}

// WriteContentType (HTML) writes HTML ContentType.
func (r HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}
