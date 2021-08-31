package main

import (
	osi "github.com/liov/hoper/server/go/lib/utils/os"
)

func main() {
	libDir, _ := osi.CMD("go list -m -f {{.Dir}}  github.com/liov/hoper/server/go/lib")
	osi.CMD("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	osi.CMD("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	osi.CMD("go install github.com/gogo/protobuf/protoc-gen-gogo")
	osi.CMD("go install " + libDir + "/tools/protoc-gen-enum")
	osi.CMD("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	osi.CMD("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	osi.CMD("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	osi.CMD("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	osi.CMD("go get -u github.com/mwitkow/go-proto-validators")
	osi.CMD("go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators")
	osi.CMD("go install " + libDir + "/tools/protoc-gen-go-patch")
}
