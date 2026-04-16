package global

type config struct {
	Locale LocaleConfig
}

type LocaleConfig struct {
	Default string
	Files map[string]string
}

func (c *config) BeforeInject() {

}

func (c *config) AfterInject() {

}
