package mail

import "gopkg.in/gomail.v2"

type mailSMTP struct {
	host     string
	port     int
	username string
	password string
}

func NewMailSMTP(host string, port int, username string, password string) mailSMTP {
	return mailSMTP{
		host:     host,
		port:     port,
		username: username,
		password: password,
	}
}

func (m mailSMTP) Send(msg ...*gomail.Message) error {
	d := gomail.NewDialer(m.host, m.port, m.username, m.password)
	return d.DialAndSend(msg...)
}
