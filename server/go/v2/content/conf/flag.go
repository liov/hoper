package conf

import "flag"

type flagValue struct {
	Password     string
	MailPassword string
}

func init() {
	flag.StringVar(&Config.Flag.Password, "p", "", "password")
	flag.StringVar(&Config.Flag.MailPassword, "mp", "", "password")
}
