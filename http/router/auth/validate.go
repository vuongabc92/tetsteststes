package auth

import (
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/session"
	"github.com/vuongabc92/octocv/http/validator"
)

func (a *authRouter) validateRegisterForm(req interface{}) *session.MessageBag {
	v := validator.New()
	return v.Validate(req, a.getRegisterValidateMessages())
}

func (a *authRouter) getRegisterValidateMessages() map[string]string {
	messages := make(map[string]string)
	messages["FullName.required"] = helpers.Trans("field_required")
	messages["FullName.alpha_space"] = helpers.Trans("unallowed_char")
	messages["FullName.min"] = helpers.Trans("fname_min")
	messages["Email.required"] = helpers.Trans("field_required")
	messages["Email.email"] = helpers.Trans("email_invalid")
	messages["Password.required"] = helpers.Trans("field_required")
	messages["Password.min"] = helpers.Trans("pass_minmax")
	messages["Password.max"] = helpers.Trans("pass_minmax")

	return messages
}

func (a *authRouter) getLoginValidateMessages() map[string]string {
	messages := make(map[string]string)
	messages["FullName.required"] = helpers.Trans("field_required")
	messages["FullName.alpha_space"] = helpers.Trans("unallowed_char")
	messages["FullName.min"] = helpers.Trans("fname_min")
	messages["Email.required"] = helpers.Trans("field_required")
	messages["Email.email"] = helpers.Trans("email_invalid")
	messages["Password.required"] = helpers.Trans("field_required")
	messages["Password.min"] = helpers.Trans("pass_minmax")
	messages["Password.max"] = helpers.Trans("pass_minmax")

	return messages
}
