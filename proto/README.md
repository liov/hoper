import "google/protobuf/empty.proto";
https://github.com/grpc-ecosystem/grpc-gateway/issues/229
解决方案
1. grpcgateway 不支持gogo，需要手动替换gw.go文件 empty.Empty->types.Empty
2. 自定义proto文件
