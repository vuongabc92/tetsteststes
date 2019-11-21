package config

import "golang.org/x/text/language"

var (
	Languages = []language.Tag{
		language.AmericanEnglish, // Fallback language
		language.Vietnamese,
	}
)
