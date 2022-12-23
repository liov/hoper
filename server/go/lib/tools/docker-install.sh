gopath=/mnt/d/SDK/gopath
codedir=/mnt/d/code/hoper/server/go/lib
goproxy=https://goproxy.io,https://goproxy.cn,direct
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install google.golang.org/protobuf/cmd/protoc-gen-go
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work/tools golang go install ./protoc-gen-enum
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install github.com/alta/protopatch/cmd/protoc-gen-go-patch
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work/tools golang go install ./protoc-gen-grpc-gin
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work golang go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work/tools golang go install ./protoc-gen-go-patch
docker run --rm -e GOPROXY=$goproxy -v $gopath:/go -v $codedir:/work -w /work/tools golang go install ./generate

cd $gopath
docker build -t jybl/golang:goprotoc -f $codedir/tools/Dockerfile .
docker push -f jybl/golang:goprotoc