package auth

import "github.com/vuongabc92/octocv"

type Backend interface {
	Register(ctx *octocv.Context, req RegisterRequest) error
	Login(ctx *octocv.Context, req LoginRequest) error
	ForgotPassword(ctx *octocv.Context, req ForgotPasswordRequest) error
}
