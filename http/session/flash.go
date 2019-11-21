package session

import (
	"encoding/json"
	"strings"
)

const flashError = "_flash_error"
const flashSuccess = "_flash_success"
const flashFormData = "_flash_form_data"

// Flash is a struct that helps with the operations over flash messages.
type Flash struct {
	session *Session
}

// Add session flash message
// https://www.gorillatoolkit.org/pkg/sessions
func (f *Flash) AddFlash(value interface{}, vars ...string) {
	f.session.Session.AddFlash(value, vars...)
	f.session.Save()
}

// Add error message in to session flash with specific flash key
// the flash message is json encoding []byte that carries signatures of MessageBag
// the error message will be remove after read.
// Check more at https://www.gorillatoolkit.org/pkg/sessions
func (f *Flash) AddError(mb *MessageBag) {
	e, _ := json.Marshal(mb)
	f.session.Session.AddFlash(e, flashError)
	f.session.Save()
}

// Return a bag of error messages that was get from session flash
// The message was encode as json []byte
// If message bag was come from form validation
// After decode json, it contains something like below:
// map[email.required:Email is required. first_name.short: First name is too short. .....]
// function GetError filters and returns a bag of messages that look like below:
// map[email:Email is required. first_name: First name is too short. .....].
// The error message will be remove after read.
// Check more at https://www.gorillatoolkit.org/pkg/sessions
func (f *Flash) GetError() *MessageBag {
	fErr := f.Flashes(flashError)
	if len(fErr) == 0 {
		return &MessageBag{}
	}

	var mb *MessageBag
	if err := json.Unmarshal(fErr[0].([]byte), &mb); err == nil {
		var uniqueMessages = NewMessageBag()
		for k, m := range mb.Messages {
			// k is a string that looks like: first_name.long, first_name.required, email.required, ...
			// what we want is: first_name, email, ...
			msgKeySplit := strings.Split(k, ".")
			if len(msgKeySplit) != 2 {
				continue
			}

			if uniqueMessages.Has(msgKeySplit[0]) {
				continue
			}

			uniqueMessages.Add(msgKeySplit[0], m)
		}

		return uniqueMessages
	}

	return &MessageBag{}
}

// Add success message in to session flash with specific flash key
// the flash message is json encoding []byte that carries signatures of MessageBag
// the success message will be remove after read.
// Check more at https://www.gorillatoolkit.org/pkg/sessions
func (f *Flash) AddSuccess(mb *MessageBag) {
	e, _ := json.Marshal(mb)
	f.session.Session.AddFlash(e, flashSuccess)
	f.session.Save()
}

// Return a bag of success messages that was get from session flash
// The message was encode as json []byte
// After decode json, it contains something like below:
// map[email.sent:Email was sent. first_name.short: First name was short. .....]
// function GetSuccess filters and returns a bag of messages that look like below:
// map[email:Email is sent. first_name: First name was short. .....].
// The success message will be remove after read.
// Check more at https://www.gorillatoolkit.org/pkg/sessions
func (f *Flash) GetSuccess() *MessageBag {
	fSuc := f.Flashes(flashSuccess)
	if len(fSuc) == 0 {
		return &MessageBag{}
	}

	var mb *MessageBag
	if err := json.Unmarshal(fSuc[0].([]byte), &mb); err == nil {
		var uniqueMessages = NewMessageBag()
		for k, m := range mb.Messages {
			// k is a string that looks like: first_name.long, first_name.required, email.required, ...
			// what we want is: first_name, email, ...
			msgKeySplit := strings.Split(k, ".")
			if len(msgKeySplit) != 2 {
				continue
			}

			if uniqueMessages.Has(msgKeySplit[0]) {
				continue
			}

			uniqueMessages.Add(msgKeySplit[0], m)
		}

		return uniqueMessages
	}

	return &MessageBag{}
}

// Add form data to flash message
// when submit form, usually data will be save here
// in case form error, the data will be read again in old form
// data will be clear after read.
func (f *Flash) AddFormData(formData *FormData) {
	d, _ := json.Marshal(formData)
	f.session.Session.AddFlash(d, flashFormData)
	f.session.Save()
}

func (f *Flash) GetFormData() *FormData {
	fErr := f.Flashes(flashFormData)
	if len(fErr) == 0 {
		return &FormData{}
	}

	var fb *FormData
	if err := json.Unmarshal(fErr[0].([]byte), &fb); err == nil {
		return fb
	}

	return &FormData{}
}

// Get all flashes by key
func (f *Flash) Flashes(vars ...string) []interface{} {
	flashes := f.session.Session.Flashes(vars...)
	f.session.Save()
	return flashes
}

//newFlash creates a new Flash and loads the session data inside its data.
func NewFlash(session *Session) *Flash {
	flash := &Flash{session: session}
	return flash
}
