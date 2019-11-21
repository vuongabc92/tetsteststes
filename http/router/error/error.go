package error

import "github.com/vuongabc92/octocv/http/router"

type authRouter struct {
	routes []router.Route
}

func NewRouter() *authRouter {
	r := &authRouter{}

	r.initRouter()

	return r
}

func (r *authRouter) Routes() []router.Route {
	return r.routes
}

func (r *authRouter) initRouter() {
	r.routes = []router.Route{}
}
