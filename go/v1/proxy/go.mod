module github.com/liov/go/proxy

go 1.13

require (
	github.com/gin-gonic/gin v1.4.0 // indirect
	github.com/liov/hoper/go/v1/proto v0.0.0-20190920064137-7615f290bf48
	google.golang.org/grpc v1.23.1
)

replace github.com/liov/hoper/go/v1/proto => ../proto
