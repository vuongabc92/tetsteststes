package middleware

import (
	"github.com/vuongabc92/octocv"
)

// Middleware is an interface to allow the use of ordinary functions as filters.
// Any struct that has the appropriate signature can be registered as a middleware.
type Middleware interface {
	WrapHandler(func(ctx *octocv.Context, vars map[string]string) error) func(ctx *octocv.Context, vars map[string]string) error
}
