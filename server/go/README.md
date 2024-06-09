# hoper3.0
cd server/go
go run $(go list -m -f {{.Dir}}  github.com/hopeio/cherry)/tools/protoc/install_tools.go
protogen.exe go -e -w -v -q -p ../../proto protogen.exe go -e -w -v -q -p ../../proto 