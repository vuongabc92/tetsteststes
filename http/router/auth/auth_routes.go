package auth

import (
	"github.com/vuongabc92/octocv"
	"net/http"
	"strings"
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

func (a *authRouter) verifyEmail(ctx *octocv.Context, vars map[string]string) error {
	var err error

	req := VerifyEmailRequest{
		Token: vars["token"],
	}

	if err = a.backend.VerifyEmail(ctx, req); err != nil {
		ctx.Logger.Errorf("Can not process to verify email address. Error: %s", err)
		return err
	}

	return nil
}

func (a *authRouter) resetPassword(ctx *octocv.Context, vars map[string]string) error {
	var err error

	req := ResetPasswordRequest{
		Token: vars["token"],
	}

	if err = a.backend.ResetPassword(ctx, req); err != nil {
		ctx.Logger.Errorf("Can not process to reset password. Error: %s", err)
		return err
	}

	return nil
}

func (a *authRouter) updatePassword(ctx *octocv.Context, vars map[string]string) error {
	var (
		err error
		req UpdatePasswordRequest
	)

	// Binding request
	if err = ctx.Bind(&req); err != nil {
		ctx.Logger.Errorf("Binding request got error: %s", err)
		return err
	}

	req.Token = vars["token"]

	// Validate form
	if msgBag := a.validateUpdatePasswordForm(req); msgBag.IsError() {
		// Add error messages to session flash
		ctx.Flash.AddError(msgBag)

		// Redirect back to register form
		ctx.Redirect(ctx.Url("get_reset_password", "token", req.Token))
		return nil
	}

	if strings.TrimSpace(req.Token) == "" {
		ctx.Redirect(ctx.Url("get_home"))
		return nil
	}

	if err = a.backend.UpdatePassword(ctx, req); err != nil {
		ctx.Logger.Errorf("Can not process to update password when user forgot password. Error: %s", err)
		return err
	}

	return nil
}
