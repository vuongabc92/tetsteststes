package mail

import "gopkg.in/gomail.v2"

type message struct {
	from    string
	to      string
	subject string
	body    string
}

func NewMessage(from string, to string, subject string, body string) message {
	return message{
		from:    from,
		to:      to,
		subject: subject,
		body:    body,
	}
}

func (m message) GetMessage() *gomail.Message {
	message := gomail.NewMessage()
	message.SetHeader("From", m.from)
	message.SetHeader("To", m.to)
	message.SetHeader("Subject", m.subject)
	message.SetBody("text/html", m.body)

	return message
}
