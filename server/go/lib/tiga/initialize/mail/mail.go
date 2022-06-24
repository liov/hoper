package mail

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

func (conf *MailConfig) generate() smtp.Auth {
	return smtp.PlainAuth("", conf.From, conf.Password, conf.Host)
}

func (conf *MailConfig) Generate() interface{} {
	return conf.generate()
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
