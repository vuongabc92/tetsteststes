package worker

import (
	"encoding/json"
	"flag"
	"github.com/go-redis/redis"
	"github.com/vuongabc92/octocv/config"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/http/mail"
	"github.com/vuongabc92/octocv/template"
	"gopkg.in/gomail.v2"
)

const (
	RegisterConfirmationChannel string = "_RegisterConfirmation"
	ForgotPasswordChannel       string = "_ForgotPassword"
)

type senderFunc func([]byte) error

type mailWorker struct {
	redis          *redis.Client
	host           string
	port           int
	username       string
	password       string
	templateEngine *template.Engine
	logger         *log.Log
}

func (m *mailWorker) run() {
	go func() {
		if err := m.sendMail(RegisterConfirmationChannel, m.sendRegisterConfirmationMail); err != nil {
			m.logger.Errorf("Send register confirmation email got error: %s", err.Error())
		}
	}()

	go func() {
		if err := m.sendMail(ForgotPasswordChannel, m.sendForgotPasswordMail); err != nil {
			m.logger.Errorf("Send forgot password email got error: %s", err.Error())
		}
	}()
}

// Waiting for signal from redis message, base on channel
// to send correct mail to user.
func (m *mailWorker) sendMail(channelName string, handler senderFunc) error {
	redisPubSub := m.redis.Subscribe(channelName)

	// Wait for confirmation that subscription is created before publishing anything.
	_, err := redisPubSub.Receive()
	if err != nil {
		return err
	}

	// Go channel which receives messages.
	messages := redisPubSub.Channel()

	// Consume messages.
	for msg := range messages {
		if err = handler([]byte(msg.Payload)); err != nil {
			m.logger.Errorf("Can not send mail for channel %s. Error: %s", channelName, err.Error())
			continue
		}
	}

	return nil
}

// Send register confirmation email over SMTP with beautiful template
func (m *mailWorker) sendRegisterConfirmationMail(msg []byte) error {
	var (
		err      error
		mailBody string
		mailData mail.RegisterConfirmationMail
	)

	flag.Parse()

	if err = json.Unmarshal(msg, &mailData); err != nil {
		m.logger.Errorf("Can not decode message was publish from register. Error: %s", err.Error())
		return err
	}

	mailBody, err = m.htmlString("mail.register-confirmation", mailData)
	from := *config.NoReplyEmailAddress
	subject := template.Trans("register_confirmation_subject")
	message := mail.NewMessage(from, mailData.MailTo, subject, mailBody)

	return m.sendMailSMTP(message.GetMessage())
}

// Send register forgot password over SMTP with amazing template :)
func (m *mailWorker) sendForgotPasswordMail(msg []byte) error {
	var (
		err      error
		mailBody string
		mailData mail.ForgotPasswordMail
	)

	flag.Parse()

	if err = json.Unmarshal(msg, &mailData); err != nil {
		m.logger.Errorf("Can not decode message was publish from forgot password. Error: %s", err.Error())
		return err
	}

	mailBody, err = m.htmlString("mail.forgot-password", mailData)
	from := *config.NoReplyEmailAddress
	subject := template.Trans("forgot_password_subject")
	message := mail.NewMessage(from, mailData.MailTo, subject, mailBody)

	return m.sendMailSMTP(message.GetMessage())
}

// Send mail over SMTP
func (m *mailWorker) sendMailSMTP(msg ...*gomail.Message) error {
	mailSmtp := mail.NewMailSMTP(m.host, m.port, m.username, m.password)

	return mailSmtp.Send(msg...)
}

// HTML renders the HTTP template specified by its file name and
// returns as string.
// See http://golang.org/doc/articles/wiki/
func (m *mailWorker) htmlString(name string, data interface{}) (string, error) {
	instance := m.templateEngine.HTMLRender.Instance(name, data)
	return instance.String()
}
