//go:build tools

package tools

import (
	_ "github.com/alta/protopatch/cmd/protoc-gen-go-patch"
	_ "github.com/googleapis/googleapis"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/mwitkow/go-proto-validators/protoc-gen-govalidators"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

////go:generate go mod download
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate protoc -I../protobuf --go_out=paths=source_relative:../protobuf ../protobuf/patch/*.proto
//go:generate protoc -I../protobuf -I../protobuf/third --go_out=paths=source_relative:../protobuf ../protobuf/utils/apiconfig/*.proto
//go:generate protoc -I../protobuf -I../protobuf/third --go_out=paths=source_relative:../protobuf ../protobuf/utils/openapiconfig/*.proto
//go:generate go install github.com/gogo/protobuf/protoc-gen-gogo
//go:generate protoc -I../protobuf -I../protobuf/utils/gogo --gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,paths=source_relative:../protobuf ../protobuf/utils/gogo/*.gen.proto
//go:generate go install ./protoc-gen-enum
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
//go:generate go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
//go:generate go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
//go:generate go install ./protoc-gen-grpc-gin
//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go get -u github.com/mwitkow/go-proto-validators
//go:generate go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
//go:generate go install ./protoc-gen-go-patch
