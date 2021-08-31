module github.com/liov/hoper/server/go/mod

go 1.17

require (
	github.com/cockroachdb/pebble v0.0.0-20210823170338-2aba043dd4a2
	github.com/dgraph-io/ristretto v0.1.0
	github.com/dgrijalva/jwt-go/v4 v4.0.0-preview1
	github.com/gin-gonic/gin v1.7.4
	github.com/go-oauth2/oauth2/v4 v4.4.1
	github.com/go-redis/redis/v8 v8.11.3
	github.com/gofiber/fiber/v2 v2.17.0
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.4.2
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.5.0
	github.com/liov/hoper/server/go/lib v1.0.0
	github.com/mwitkow/go-proto-validators v0.3.2
	github.com/spf13/viper v1.8.1
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.19.0
	golang.org/x/net v0.0.0-20210813160813-60bc85c4be6d
	google.golang.org/genproto v0.0.0-20210821163610-241b8fcbd6c8
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/gorm v1.21.13
)

replace github.com/liov/hoper/server/go/lib => ../lib
