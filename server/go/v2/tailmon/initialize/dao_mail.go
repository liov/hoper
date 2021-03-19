package initialize

import (
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/reflect"
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

func (init *Init) P3Mail() smtp.Auth {
	conf := &MailConfig{}
	if exist := reflecti.GetFieldValue(init.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}
