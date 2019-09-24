module github.com/liov/hoper/go/v1/helloworld

go 1.13

require (
	github.com/liov/hoper/go/v1/initialize v0.0.0-20190920064137-7615f290bf48 // indirect
	github.com/liov/hoper/go/v1/protobuf v0.0.0-20190920064137-7615f290bf48 // indirect
	github.com/liov/hoper/go/v1/utils v0.0.0-20190920064137-7615f290bf48
	google.golang.org/grpc v1.23.1
)

replace (
	github.com/liov/hoper/go/v1/initialize => ../initialize
	github.com/liov/hoper/go/v1/protobuf => ../protobuf
	github.com/liov/hoper/go/v1/utils => ../utils
)
