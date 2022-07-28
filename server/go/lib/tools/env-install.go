package main

import (
	osi "github.com/actliboy/hoper/server/go/lib/utils/os"
)

// 提供给使用框架的人安装所需环境
func main() {
	libDir, _ := osi.CMD("go list -m -f {{.Dir}}  github.com/actliboy/hoper/server/go/lib")
	osi.CMD("go install google.golang.org/protobuf/cmd/protoc-gen-go")
	/*	osi.CMD("protoc -I" + libDir + "/protobuf -I" + libDir + "/protobuf/third --go_out=paths=source_relative:" + libDir + "/protobuf" + libDir + "/protobuf/utils/patch/*.proto")
		osi.CMD("protoc -I" + libDir + "/protobuf -I" + libDir + "/protobuf/third --go_out=paths=source_relative:" + libDir + "/protobuf" + libDir + "/protobuf/utils/apiconfig/*.proto")
		osi.CMD("protoc -I" + libDir + "/protobuf -I" + libDir + "/protobuf/third --go_out=paths=source_relative:" + libDir + "/protobuf" + libDir + "/protobuf/utils/openapiconfig/*.proto")*/
	osi.CMD("go install " + libDir + "/tools/protoc-gen-grpc-gin")
	osi.CMD("go install " + libDir + "/tools/protoc-gen-enum")
	osi.CMD("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway")
	osi.CMD("go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2")
	osi.CMD("go install github.com/alta/protopatch/cmd/protoc-gen-go-patch")
	osi.CMD("go install google.golang.org/grpc/cmd/protoc-gen-go-grpc")
	osi.CMD("go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators")
	osi.CMD("go install " + libDir + "/tools/protoc-gen-go-patch")
	osi.CMD("go install " + libDir + "/generate")
}
