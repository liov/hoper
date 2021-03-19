//+build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/liov/hoper/go/v2/tools/protoc-gen-go-patch"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go:generate go install ./protoc-gen-go-enum
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
//go:generate go install github.com/liov/hoper/go/v2/tools/protoc-gen-go-patch
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
