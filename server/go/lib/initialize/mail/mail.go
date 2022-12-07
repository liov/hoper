package mail

import (
	"github.com/liov/hoper/server/go/lib/initialize"
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

func (conf *MailConfig) Build() smtp.Auth {
	return smtp.PlainAuth("", conf.From, conf.Password, conf.Host)
}

func (conf *MailConfig) Generate() interface{} {
	return conf.Build()
}

type Mail struct {
	smtp.Auth
	Conf MailConfig
}

func (m *Mail) Config() initialize.Generate {
	return &m.Conf
}

func (m *Mail) SetEntity(entity interface{}) {
	if client, ok := entity.(smtp.Auth); ok {
		m.Auth = client
	}
}
