module github.com/liov/hoper/go/v1/initialize

go 1.13

require (
	github.com/bluele/gcache v0.0.0-20190518031135-bc40bd653833
	github.com/etcd-io/bbolt v1.3.3
	github.com/garyburd/redigo v1.6.0
	github.com/globalsign/mgo v0.0.0-20181015135952-eeefdecb41b8
	github.com/jinzhu/configor v1.1.1
	github.com/jinzhu/gorm v1.9.10
	github.com/liov/hoper/go/v1/utils v0.0.0-20190920064137-7615f290bf48
	go.uber.org/zap v1.10.0
)

replace github.com/liov/hoper/go/v1/utils => ../utils
