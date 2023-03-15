package mail

import (
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

type Mail struct {
	smtp.Auth
	Conf MailConfig
}

func (m *Mail) Config() any {
	return &m.Conf
}

func (m *Mail) SetEntity() {
	m.Auth = m.Conf.Build()
}

func (m *Mail) Close() error {
	return nil
}
