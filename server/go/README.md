# hoper3.0
cd server/go
go run $(go list -m -f {{.Dir}}  github.com/hopeio/protobuf)/tools/install_tools.go
windows:`protogen.exe go -d -e -w -v -i ../../proto`
unix:`protogen go -d -e -w -v -i ../../proto`

windows:`protogen.exe go -d -e -w -v -i ../../proto  (graphql)`
unix:`protogen go -d -e -w -v -i ../../proto  (graphql)`
# docker
GOOS=linux go build -tags go-json -o ../../build/hoper