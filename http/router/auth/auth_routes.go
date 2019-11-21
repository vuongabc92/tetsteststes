package auth

import (
	"github.com/vuongabc92/octocv"
	"net/http"
)

func (*authRouter) loadLoginPage(ctx *octocv.Context, vars map[string]string) error {
	ctx.HTML(http.StatusOK, "auth.login", nil)

	return nil
}

func (*authRouter) loadRegisterPage(ctx *octocv.Context, vars map[string]string) error {
	ctx.HTML(http.StatusOK, "auth.register", nil)

	return nil
}

func (*authRouter) loadForgotPasswordPage(ctx *octocv.Context, vars map[string]string) error {
	ctx.HTML(http.StatusOK, "auth.forgot-password", nil)

	return nil
}

func (a *authRouter) login(ctx *octocv.Context, vars map[string]string) error {
	var (
		err error
		req LoginRequest
	)

	// Binding request
	if err = ctx.Bind(&req); err != nil {
		ctx.Logger.Errorf("Binding request got error: %s", err)
		return err
	}

	if err = a.backend.Login(ctx, req); err != nil {
		ctx.Logger.Errorf("Login got error: %s", err)
		return err
	}

	return nil
}

func (a *authRouter) register(ctx *octocv.Context, vars map[string]string) error {
	var (
		err error
		req RegisterRequest
	)

	// Binding request
	if err = ctx.Bind(&req); err != nil {
		ctx.Logger.Errorf("Binding request got error: %s", err)
		return err
	}

	// Validate form
	if msgBag := a.validateRegisterForm(req); msgBag.IsError() {
		// Add error messages to session flash
		ctx.Flash.AddError(msgBag)

		// Redirect back to register form
		ctx.Redirect(ctx.Url("get_register"))
		return nil
	}

	if err = a.backend.Register(ctx, req); err != nil {
		ctx.Logger.Errorf("Can not register user, email: %s. Error: %s", req.Email, err.Error())
		return err
	}

	return nil
}

func (a *authRouter) forgotPassword(ctx *octocv.Context, vars map[string]string) error {
	var (
		err error
		req ForgotPasswordRequest
	)

	// Binding request
	if err = ctx.Bind(&req); err != nil {
		ctx.Logger.Errorf("Binding request got error: %s", err)
		return err
	}

	if err = a.backend.ForgotPassword(ctx, req); err != nil {
		ctx.Logger.Errorf("Can not process to send reclaim password instruction. Error: %s", err)
		return err
	}

	return nil
}
