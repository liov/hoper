package initialize

import (
	"net/smtp"
)

type MailConfig struct {
	Host     string
	Port     string
	From     string
	Password string
}

func (conf *MailConfig) Generate() smtp.Auth {
	return smtp.PlainAuth("", conf.From, conf.Password, conf.Host)
}

func (init *Inject) P3Mail() smtp.Auth {
	conf := &MailConfig{}
	if exist := init.SetConf(conf); !exist {
		return nil
	}
	return conf.Generate()
}
