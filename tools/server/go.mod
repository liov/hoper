module tools

go 1.15

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/liov/hoper/server/go/lib v1.0.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/tidwall/gjson v1.9.3
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
	gorm.io/gorm v1.21.13
)

replace github.com/liov/hoper/server/go/lib => ../../server/go/lib
