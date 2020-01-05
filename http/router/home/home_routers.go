package home

import (
	"github.com/vuongabc92/octocv"
	"net/http"
)

func (*homeRouter) loadHomePage(ctx *octocv.Context, vars map[string]string) error {
	ctx.HTML(http.StatusOK, "page.home", nil)

	return nil
}
