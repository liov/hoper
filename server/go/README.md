# hoper3.0
cd server/go
go run $(go list -m -f {{.Dir}}  github.com/hopeio/protobuf)/tools/install_tools.go
windows:`protogen.exe go -e -w -v -p ../../proto`
unix:`protogen go -e -w -v -p ../../proto`

windows:`protogen.exe go -e -w -v -q -p ../../proto  (graphql)`
unix:`protogen go -e -w -v -q -p ../../proto  (graphql)`
# docker
GOOS=linux go build -tags go-json -o ../../build/hoper