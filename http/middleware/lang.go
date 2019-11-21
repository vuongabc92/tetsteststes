package middleware

import (
	"github.com/vuongabc92/octocv"
	"github.com/vuongabc92/octocv/config"
	octolang "github.com/vuongabc92/octocv/lang"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Lang struct {
}

func NewLangMiddleware() Lang {
	return Lang{}
}

func (l Lang) WrapHandler(handler func(ctx *octocv.Context, vars map[string]string) error) func(ctx *octocv.Context, vars map[string]string) error {
	return func(ctx *octocv.Context, vars map[string]string) error {
		matcher := language.NewMatcher(config.Languages)
		lang, _ := ctx.Request.Cookie("lang")
		accept := ctx.Request.Header.Get("Accept-Language")
		tag, _ := language.MatchStrings(matcher, lang.String(), accept)

		// Set language tag to get the correct texts
		octolang.Trans = message.NewPrinter(tag)

		return handler(ctx, vars)
	}
}
