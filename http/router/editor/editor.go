package editor

import (
	"github.com/vuongabc92/octocv/http/router"
)

type editorRouter struct {
	backend Backend
	routes  []router.Route
}

func NewRouter(backend Backend) *editorRouter {
	r := &editorRouter{
		backend: backend,
	}

	r.initRouter()

	return r
}

func (r *editorRouter) Routes() []router.Route {
	return r.routes
}

func (r *editorRouter) initRouter() {
	r.routes = []router.Route{
		//Get
		router.NewGetRoute("editor", r.ShowEditorPage),
	}
}
