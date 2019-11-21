package validator

import (
	"fmt"
	"github.com/vuongabc92/octocv/helpers"
	"github.com/vuongabc92/octocv/http/session"
	"gopkg.in/go-playground/validator.v9"
	"regexp"
)

const customRuleAlphaSpace = "alpha_space"

type Validator struct {
}

func New() *Validator {
	return &Validator{}
}

func (v *Validator) Validate(req interface{}, msgBag map[string]string) *session.MessageBag {
	var (
		msg      = session.NewMessageBag()
		validate = validator.New()
	)

	//Register custom validation rule
	if err := validate.RegisterValidation(customRuleAlphaSpace, v.alphaSpace); err != nil {
		panic("Can not register custom rule. Error: " + err.Error())
	}

	if err := validate.Struct(req); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			key := fmt.Sprintf("%s.%s", e.Field(), e.Tag())
			if m, ok := msgBag[key]; ok {
				msg.Add(helpers.ToSnakeCase(key), m)
			}
		}
	}

	return msg
}

// Validate that field value must contain alpha and space only
func (v *Validator) alphaSpace(fl validator.FieldLevel) bool {
	pattern := `^[a-zA-Z\s]+$`
	return regexp.MustCompile(pattern).MatchString(fl.Field().String())
}
