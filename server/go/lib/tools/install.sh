go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc -I../protobuf -I../protobuf/third --go_out=paths=source_relative:../protobuf ../protobuf/utils/**/*.proto
go install ./protoc-gen-enum
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
go install ./protoc-gen-grpc-gin
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
go install ./protoc-gen-go-patch
go install ./generate
