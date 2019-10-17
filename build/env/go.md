sudo tar -C /usr/local -xzf go1.13.1.linux-amd64.tar.gz
vi /etc/
export PATH=$PATH:/usr/local/go/bin
export GOROOT=$PATH:/usr/local/go
export GOPATH=/mnt/e/gopath

cmd
set GOARCH=amd64
set GOOS=linux
go build

ps
$env:GOARCH="amd64"
$env:GOOS="linux"
go build

gitbash

export GOARCH=amd64
export GOOS=linux
go build

hosts 192.30.253.112	github.com
      199.232.5.194	github.global.ssl.fastly.net
go get -u github.com/gpmgo/gopm
安装protoc[https://github.com/protocolbuffers/protobuf/releases]
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/micro/micro
go get -u github.com/micro/protoc-gen-micro