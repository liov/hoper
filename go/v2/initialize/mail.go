package initialize

import (
	"net/smtp"

	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P3Mail() smtp.Auth {
	conf := MailConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &conf); !exist {
		return nil
	}
	return smtp.PlainAuth("", conf.User, conf.Password, conf.Host)
}
