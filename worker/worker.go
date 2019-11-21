package worker

import (
	"github.com/go-redis/redis"
	"github.com/vuongabc92/octocv/http/log"
	"github.com/vuongabc92/octocv/http/mail"
	"github.com/vuongabc92/octocv/template"
)

type Option struct {
	Redis          *redis.Client
	SMTPMail       *mail.SMTPMail
	TemplateEngine *template.Engine
	Logger         *log.Log
}

type worker struct {
	Options *Option
	Mail    *mailWorker
}

func NewWorker(opt *Option) worker {
	mail := &mailWorker{
		redis:          opt.Redis,
		host:           opt.SMTPMail.Host,
		port:           opt.SMTPMail.Port,
		username:       opt.SMTPMail.Username,
		password:       opt.SMTPMail.Password,
		templateEngine: opt.TemplateEngine,
		logger:         opt.Logger,
	}

	return worker{
		Options: opt,
		Mail:    mail,
	}
}

func (w *worker) Run() {
	go w.Mail.run()
}
