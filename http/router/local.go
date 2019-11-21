package router

import (
	"github.com/vuongabc92/octocv/http/httputils"
)

//Local route implements router.Route
type localRoute struct {
	method  string
	path    string
	handler httputils.HandlerFunc
	name    string
}

func (r localRoute) Handler() httputils.HandlerFunc {
	return r.handler
}

func (r localRoute) Method() string {
	return r.method
}

func (r localRoute) Path() string {
	return r.path
}

func (r localRoute) Name() string {
	return r.name
}

func NewRoute(method, path string, handler httputils.HandlerFunc, name string) Route {
	return localRoute{method, path, handler, name}
}

func NewGetRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("GET", path, handler, name)
}

func NewPostRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("POST", path, handler, name)
}

func NewPutRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("PUT", path, handler, name)
}

func NewDeleteRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("DELETE", path, handler, name)
}

func NewOptionsRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("OPTIONS", path, handler, name)
}

func NewHeadRoute(path string, handler httputils.HandlerFunc, name string) Route {
	return NewRoute("HEAD", path, handler, name)
}
