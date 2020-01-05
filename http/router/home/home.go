package home

import (
	"github.com/vuongabc92/octocv/http/router"
)

type homeRouter struct {
	backend Backend
	routes  []router.Route
}

func NewRouter(backend Backend) *homeRouter {
	r := &homeRouter{
		backend: backend,
	}

	r.initRouter()

	return r
}

func (h *homeRouter) Routes() []router.Route {
	return h.routes
}

func (h *homeRouter) initRouter() {
	h.routes = []router.Route{
		//Get
		router.NewGetRoute("/", h.loadHomePage, "get_home"),
	}
}
