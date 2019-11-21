package lang

import (
	"github.com/vuongabc92/octocv/lang/locales"
	"golang.org/x/text/language"
)

// list texts per language is supported
var texts = map[language.Tag]Text{
	language.AmericanEnglish: locales.EnUsTexts,
	language.Vietnamese:      locales.ViTexts,
}

type Text map[string]string

func NewText() Text {
	return Text{}
}

func (l Text) GetTexts(tag language.Tag) map[string]string {
	return texts[tag]
}
