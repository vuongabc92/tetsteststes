package auth

import (
	"github.com/vuongabc92/octocv/http/router"
)

type authRouter struct {
	backend Backend
	routes  []router.Route
}

func NewRouter(backend Backend) *authRouter {
	r := &authRouter{
		backend: backend,
	}

	r.initRouter()

	return r
}

func (a *authRouter) Routes() []router.Route {
	return a.routes
}

func (a *authRouter) initRouter() {
	a.routes = []router.Route{
		//GET
		router.NewGetRoute("/login", a.loadLoginPage, "get_login"),
		router.NewGetRoute("/register", a.loadRegisterPage, "get_register"),
		router.NewGetRoute("/forgot-password", a.loadForgotPasswordPage, "get_forgot_password"),
		router.NewGetRoute("/verify/{token}", a.verifyEmail, "get_verify_email"),
		router.NewGetRoute("/reset/{token}", a.resetPassword, "get_reset_password"),

		//POST
		router.NewPostRoute("/login", a.login, "post_login"),
		router.NewPostRoute("/register", a.register, "post_register"),
		router.NewPostRoute("/forgot-password", a.forgotPassword, "post_forgot_password"),
		router.NewPostRoute("/update-password/{token}", a.updatePassword, "post_update_password"),
	}
}
