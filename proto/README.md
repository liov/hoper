import "google/protobuf/empty.proto";
https://github.com/grpc-ecosystem/grpc-gateway/issues/229
解决方案
1. grpcgateway 不支持gogo，需要手动替换gw.go文件 empty.Empty->types.Empty
2. 自定义proto文件

必须要有一个规范，此规范为了分辨及用于生成普通proto
凡仅做引用：
1.因package问题的，实际用于生成的proto文件命名为*.gen.proto,引用的文件命名为*.imp.proto
且*.gen.proto文件中不应出现option go_package
~~2.因不对外开放的enum，引用的proto文件命名为*.enum.pub.proto~~ 因为一些原因，这条去掉
这些枚举值秘密程度没有那么高，且不一定会对外暴露，若要避免对外暴露可以不暴露enum.proto，让使用方自行定义

这么做的原因是包名定义github.com/xxx,然而项目又不在github.com目录下，
在指定pb.go的生成目录(非github.com/xxx)，这样就会造成生成目录下自动创建github.com/xxx，造成错位

E:\protoc\bin\protoc.exe -ID:\hoper\server\go/../../proto
-IE:\gopath\pkg\mod\github.com\grpc-ecosystem\grpc-gateway\v2@v2.1.0
-IE:\gopath\pkg\mod\github.com\grpc-ecosystem\grpc-gateway\v2@v2.1.0/third_party/googleapis
-IE:\gopath\pkg\mod\google.golang.org\protobuf@v1.26.0  
-ID:\hoper\server\go/../../proto/utils/proto
-IE:\gopath/src D:\hoper\server\go/../../proto/user/*model.proto
--dart_out=plugin=go,paths=source_relative:D:\hoper\client\flutter\generated\protobuf
