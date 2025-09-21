cd $(dirname $0) && pwd
protogen go -e -w -v -p ../../proto -o protobuf
GOOS=linux go build -trimpath -o ../../build/hoper main.go
cd ../../build
docker build -t jybl/hoper . && docker push jybl/hoper

