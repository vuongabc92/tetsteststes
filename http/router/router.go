package router

import (
	"github.com/vuongabc92/octocv/http/httputils"
)

type Router interface {
	Routes() []Route
}

type Route interface {
	Handler() httputils.HandlerFunc
	Method() string
	Path() string
	Name() string
}
