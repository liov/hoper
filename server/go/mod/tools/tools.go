//+build tools

package tools

import (
	_ "github.com/alta/protopatch/cmd/protoc-gen-go-patch"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go list -m -f {{.Dir}} github.com/liov/hoper/server/go/lib
//go:generate go run D:/hoper/server/go/lib/tools/install.go
