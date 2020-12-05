module tools

go 1.15

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/kataras/iris/v12 v12.0.1
	github.com/liov/hoper/go/v2 v2.0.0
	github.com/mozillazg/go-pinyin v0.18.0
	github.com/robfig/cron v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/tidwall/gjson v1.6.3
	golang.org/x/net v0.0.0-20201002202402-0a1ea396d57c
	gorm.io/gorm v1.20.7
)

replace github.com/liov/hoper/go/v2 => ../../server/go/v2
