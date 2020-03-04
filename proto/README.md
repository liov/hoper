import "google/protobuf/empty.proto";
https://github.com/grpc-ecosystem/grpc-gateway/issues/229
解决方案
1. grpcgateway 不支持gogo，需要手动替换gw.go文件 empty.Empty->types.Empty
2. 自定义proto文件

必须要有一个规范，此规范为了分辨及用于生成普通proto
凡仅做引用：
1.因package问题的，实际用于生成的proto文件命名为*.gen.proto,引用的文件命名为*.imp.proto
且*.gen.proto文件中不应出现option go_package
2.因不对外开放的enum，引用的proto文件命名为*.enum.pub.proto