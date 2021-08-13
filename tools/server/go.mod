module tools

go 1.15

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/liov/hoper/v2 v2.0.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/tidwall/gjson v1.6.3
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	gorm.io/gorm v1.21.8
)

replace github.com/liov/hoper/v2 => ../../server/go
