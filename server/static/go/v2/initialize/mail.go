package initialize

import (
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/reflect3"
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

func (i *Init) P3Mail() smtp.Auth {
	conf := &MailConfig{}
	if exist := reflect3.GetFieldValue(i.conf, conf); !exist {
		return nil
	}
	return conf.Generate()
}
