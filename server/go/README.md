# hoper3.0
cd server/go
go run $(go list -m -f {{.Dir}}  github.com/hopeio/cherry)/tools/protoc/install_tools.go
protogen.exe go -e -w -v -p ../../proto

protogen go -e -w -v -p ../../proto

protogen.exe go -e -w -v -q -p ../../proto  (graphql)
protogen go -e -w -v -q -p ../../proto  (graphql)
# docker
GOOS=linux go build -tags go-json -o ../../build/hoper