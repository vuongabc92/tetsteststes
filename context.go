package octocv

import (
	"context"
	"github.com/go-redis/redis"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	database "github.com/vuongabc92/octocv/database/mongodb"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/binding"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/http/session"
	"github.com/vuongabc92/octocv/render"
	"github.com/vuongabc92/octocv/repositories"
	"github.com/vuongabc92/octocv/template"
	gotemplate "html/template"
	"net/http"
)

type ViewVars struct {
	CsrfField gotemplate.HTML
}

type ViewData struct {
	Vars      ViewVars
	Errors    *session.MessageBag
	Successes *session.MessageBag
	FormData  *session.FormData
	Resp      interface{}
}

type Context struct {
	context.Context
	Writer         http.ResponseWriter
	Request        *http.Request
	TemplateEngine *template.Engine
	Router         *mux.Router
	Session        *session.Session
	Flash          *session.Flash
	DB             *database.MongoDBConnection
	Logger         *log.Log
	Redis          *redis.Client
}

// ShouldBind checks the Content-Type to select a binding engine automatically,
// Depending the "Content-Type" header different bindings are used:
//     "application/json" --> JSON binding
//     "application/xml"  --> XML binding
// otherwise --> returns an error
// It parses the request's body as JSON if Content-Type == "application/json" using JSON or XML as a JSON input.
// It decodes the json payload into the struct specified as a pointer.
// Like c.Bind() but this method does not set the response status code to 400 and abort if the json is not valid.
func (c *Context) Bind(obj interface{}) error {
	b := binding.Default(c.Request.Method, c.ContentType())
	return c.ShouldBindWith(obj, b)
}

// ShouldBindWith binds the passed struct pointer using the specified binding engine.
// See the binding package.
func (c *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(c.Request, obj)
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return helpers.FilterFlags(c.requestHeader("Content-Type"))
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeader(code)
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

// HTML renders the HTTP template specified by its file name.
// It also updates the HTTP code and sets the Content-Type as "text/html".
// See http://golang.org/doc/articles/wiki/
func (c *Context) HTML(code int, name string, data interface{}) {

	vars := ViewVars{
		CsrfField: csrf.TemplateField(c.Request),
	}

	viewData := ViewData{
		Vars:      vars,
		Errors:    c.Flash.GetError(),
		Successes: c.Flash.GetSuccess(),
		FormData:  c.Flash.GetFormData(),
		Resp:      data,
	}

	instance := c.TemplateEngine.HTMLRender.Instance(name, viewData)

	c.Render(code, instance)
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) Redirect(location string) {
	c.RedirectWithStatus(http.StatusFound, location)
}

// Redirect returns a HTTP redirect to the specific location.
func (c *Context) RedirectWithStatus(code int, location string) {
	c.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  c.Request,
	})
}

// Get repository factory
func (c *Context) GetRepositoryFactory() *repositories.RepositoryFactory {
	repoFactory := repositories.NewRepositoryFactory(c.DB.Database())
	return repoFactory
}

func (c *Context) Url(routeName string, queryPairs...string) string {
	url, err := c.Router.Get(routeName).URL(queryPairs...)
	//r := c.Router.Host("").GetName()
	if err != nil {
		c.Logger.Error("Get route url error: " + err.Error())
		return ""
	}

	return url.EscapedPath()
}
