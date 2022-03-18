module tools

go 1.15

require (
	github.com/PuerkitoBio/goquery v1.6.0
	github.com/actliboy/hoper/server/go/lib v1.0.0
	github.com/go-redis/redis/v8 v8.11.4 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/tidwall/gjson v1.9.3
	github.com/valyala/fasthttp v1.34.0 // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f
	gorm.io/gorm v1.21.13
)

replace github.com/actliboy/hoper/server/go/lib => ../../server/go/lib
