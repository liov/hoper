package config

import "flag"

type flagValue struct {
	Password string
	MailPassword string
	Additional string
}

func init()  {
	flag.StringVar(&Conf.Flag.Password, "p", "", "password")
	flag.StringVar(&Conf.Flag.MailPassword, "mp", "", "password")
	flag.StringVar(&Conf.Flag.Additional,"a","","额外的配置")
}
