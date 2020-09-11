//+build tools

package tools

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/liov/protopatch/cmd/protoc-gen-go-patch"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)

//go:generate go install ./protoc-gen-go-enum
