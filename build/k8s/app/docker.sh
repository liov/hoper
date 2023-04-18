GOPATH=/mnt/d/SDK/gopath
Code=/mnt/d/code/hoper
GOPROXY=https://goproxy.io
GOIMAGE=golang:1.20

# go mod tidy
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/server/go/lib -e GOPROXY=$GOPROXY $GOIMAGE go mod tidy

# goprotoc generate
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/lib/tools/generate -e GOPROXY=$GOPROXY $GOIMAGE generate go --proto=/work/proto --genpath=/work/server/go/mod/protobuf

# node
docker run --rm --privileged=true -v /home/ghoper:/work -w /work/website/vhoper node:16-alpine3.16 npm run build
docker run -v /home/ghoper/static:/static --net=host --restart=always --cpus=0.2 -d --name vhoper  vhoper:1.2