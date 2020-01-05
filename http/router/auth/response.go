package auth

import "html/template"

type RegisterResponse struct {
	CSRFField template.HTML
}

type ResetPasswordResponse struct {
	Token string
}
