cd $(dirname $0) && pwd
protogen go -d -e -w -v -i ../../proto -o protobuf
GOOS=linux go build -trimpath -o ../../build/hoper main.go
cd ../../build
docker build -t jybl/hoper . && docker push jybl/hoper

