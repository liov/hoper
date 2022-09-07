GOPATH=/mnt/d/SDK/gopath
Code=/mnt/d/code/hoper
GOPROXY=https://goproxy.io
GOIMAGE=golang:1.19

# go mod tidy
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/server/go/lib -e GOPROXY=$GOPROXY $GOIMAGE go mod tidy

# install tools
docker run --rm -v $GOPATH:/go -v $Code/server/go/:/work -w /work/lib/tools -e GOPROXY=$GOPROXY $GOIMAGE bash tools.sh

# install generate
docker run --rm -v $GOPATH:/go -v $Code/server/go/:/work -w /work/lib/tools/generate -e GOPROXY=$GOPROXY $GOIMAGE go install

# build goprotoc
docker build  --build-arg gopath=$GOPATH -t jybl/goprotoc -f $Code/server/go/lib/tools/Dockerfile  $GOPATH/bin
docker push jybl/goprotoc

# build notify
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/server/go/lib -e GOPROXY=$GOPROXY $GOIMAGE go build -o /work/build/tmp/notify ./tools/notify
docker build -t jybl/notify build/tmp/
docker push jybl/notify

# goprotoc generate
docker run --rm -v $GOPATH:/go -v $Code:/work -w /work/lib/tools/generate -e GOPROXY=$GOPROXY $GOIMAGE generate go --proto=/work/proto --genpath=/work/server/go/mod/protobuf

# test notify
docker run --rm -e PLUGIN_DING_TOKEN=xxx -e PLUGIN_DING_SECRET=xxx jybl/notify

