# hoper2.0

## quickstart

```sh
cd tools && go generate
```
user.model.proto
```protobuf
syntax = "proto3";
package user;
import "user/user.enum.proto";
import "patch/go.proto";
import "protoc-gen-openapiv2/options/annotations.proto";

option java_package = "xyz.hoper.protobuf.user";
option go_package = "github.com/actliboy/hoper/server/go/lib/protobuf/user";
// 用户
message User {
  uint64 id = 1 [(go.field) = {tags:'gorm:"primaryKey;"'}];
  string name = 2 [(go.field) = {tags:'gorm:"size:10;not null" annotation:"昵称"'}];
    // 性别，0未填写，1男，2女
  Gender gender = 8 [(go.field) = {tags:'gorm:"type:int2;default:0"'}, (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = {
    type:INTEGER
  }];
}
```
user.enum.proto
```protobuf
syntax = "proto3";
package user;
import "utils/proto/gogo/enum.proto";
import "patch/go.proto";

option (gogo.enum_gqlgen_all) = true;

option java_package = "xyz.hoper.protobuf.user";
option go_package = "github.com/actliboy/hoper/server/go/lib/protobuf/user";

option (gogo.enum_prefix_all) = false;
option (go.file) = {no_enum_prefix:true};
// 用户性别
enum Gender{
    option (go.enum) = {stringer_name: 'OrigString'};
    GenderPlaceholder = 0 [(gogo.enumvalue_cn)= "占位"];
    GenderUnfilled = 1 [(gogo.enumvalue_cn)= "未填"];
    GenderMale = 2 [(gogo.enumvalue_cn)= "男"];
    GenderFemale = 3 [(gogo.enumvalue_cn)= "女"];
}

```

user.service.proto
```protobuf
syntax = "proto3";
package user;
import "user/user.model.proto";
import "user/user.enum.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "github.com/mwitkow/go-proto-validators/validator.proto";
import "google/api/annotations.proto";
import "utils/empty/empty.proto";
import "utils/response/response.proto";
import "utils/request/param.proto";
import "utils/proto/gogo/graphql.proto";
import "utils/oauth/oauth.proto";
import "patch/go.proto";
import "google/protobuf/wrappers.proto";
option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_swagger) = {
  info: {
    version: "1.0"
  }
};

service UserService {

  option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_tag) = {
    description: "用户相关接口"
  };
    //获取用户信息
  rpc Info (request.Object) returns (User) {
    option (google.api.http) = {
      get: "/api/v1/user/{id}"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      tags:["用户相关接口", "v1.0.0"]
      summary : "获取用户信息"
      description : "根据Id获取用户信息接口"
    };
    option (gogo.graphql_operation) = Query;
  }

}
```
```sh
cd protobuf && go generate
```
business mod
```go
type config struct {
	//自定义的配置
	Customize serverConfig
	Server    initialize.ServerConfig
	Mail      initialize.MailConfig
	GORMDB    initialize.DatabaseConfig
	Redis     initialize.RedisConfig
	Cache     initialize.CacheConfig
	Log       initialize.LogConfig
	Viper     *viper.Viper
}

var Conf = &config{}

func (c *config) Init() {
	if runtime.GOOS == "windows" {
	}

	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
type serverConfig struct {
	Volume fs.Dir

	PassSalt string
	// 天数
	TokenMaxAge time.Duration
	TokenSecret string
	PageSize    int8

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string
}

// dao dao.
type dao struct {
	// GORMDB 数据库连接
	GORMDB   *gorm.DB `config:"database"`
	StdDB    *sql.DB
	PebbleDB *pebble.DB
	// RedisPool Redis连接池
	Redis *redis.Client
	Cache *ristretto.Cache
	//elastic
	MailAuth smtp.Auth
}

// CloseDao close the resource.
func (d *dao) Close() {
	if d.PebbleDB != nil {
		d.PebbleDB.Close()
	}
	if d.Redis != nil {
		d.Redis.Close()
	}
	if d.GORMDB != nil {
		rawDB, _ := d.GORMDB.DB()
		rawDB.Close()
	}
}

func (d *dao) Init() {
	db := d.GORMDB
	db.Callback().Create().Remove("gorm:save_before_associations")
	db.Callback().Create().Remove("gorm:save_after_associations")
	db.Callback().Update().Remove("gorm:save_before_associations")
	db.Callback().Update().Remove("gorm:save_after_associations")

	d.StdDB, _ = db.DB()
}
```
main.go
```go
package main

import (
	"github.com/actliboy/hoper/server/go/lib/tailmon/pick"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/gin/handler"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/actliboy/hoper/server/go/lib/protobuf/user"
	"github.com/actliboy/hoper/server/go/lib/tailmon"
	"github.com/actliboy/hoper/server/go/lib/tailmon/initialize"
	uconf "github.com/actliboy/hoper/server/go/lib/user/conf"
	udao "github.com/actliboy/hoper/server/go/lib/user/dao"
	userservice "github.com/actliboy/hoper/server/go/lib/user/service"
	"github.com/actliboy/hoper/server/go/lib/utils/log"

	"go.opencensus.io/examples/exporter"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(uconf.Conf, udao.Dao)()
	view.RegisterExporter(&exporter.PrintExporter{})
	view.SetReportingPeriod(time.Second)
	// Register the view to collect gRPC client stats.
	if err := view.Register(ocgrpc.DefaultClientViews...); err != nil {
		log.Fatal(err)
	}
	pick.RegisterService(userservice.GetUserService())
	(&tailmon.Server{
		//为了可以自定义中间件
		GRPCOptions: []grpc.ServerOption{
			grpc.ChainUnaryInterceptor(),
			grpc.ChainStreamInterceptor(),
			//grpc.StatsHandler(&ocgrpc.ServerHandler{})
		},
		GRPCHandle: func(gs *grpc.Server) {
			user.RegisterUserServiceServer(gs, userservice.GetUserService())
		},
		GinHandle: func(app *gin.Engine) {
			_ = user.RegisterUserServiceHandlerServer(app, userservice.GetUserService())
			app.Static("/static", "F:/upload")
			pick.Gin(app, true, initialize.InitConfig.Module)
		},
	}).Start()
}

```

## 目录结构
 ```sh
     - protobuf
        - generate
     - tailmon
        - context
        - intialize
        - pick
     - tools
     - utils
     - business mod
 ```

### proto
 以protobuf为核心,编写proto文件，生成所需的模型及服务  
 proto文件应该是共享的，所有模块可见的，以实现模块间rpc调用无需重复变更模型  
 本项目为实现跨语言的grpc，将proto文件抽离单独放置一个文件夹内

#### generate
 为go语言写的protobuf生成程序

### tools
本项目需要用到的protobuf插件，运行tools.go文件的//go:generate，会自动安装
- google.golang.org/protobuf/cmd/protoc-gen-go
- google.golang.org/grpc/cmd/protoc-gen-go-grpc
- github.com/gogo/protobuf/protoc-gen-gogo
- ./protoc-gen-enum 自己写的enum生成插件，分为错误enum及普通enum，生成性能更高的`String()`,错误enum会额外生成`Error()string`
- github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 用于生成swagger文档
- protoc-gen-go-patch 在github.com/alta/protopatch/cmd/protoc-gen-go-patch（利用修改ast重新生成源码）基础上，提供自定义json，去enum prefix等功能
- ./protoc-gen-grpc-gin github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway的gin版本，用于生成http接口
- github.com/mwitkow/go-proto-validators/protoc-gen-govalidators 结构体校验，利用生成代码而不是反射，性能更好


### tailmon
基于反射自动注入的自研框架，整合grpc及http调用，使proto定义的服务可对外分别提供grpc，http及可选的graphql接口
集成opencensus实现调用链路跟踪记录，配合context及utils-log 实现完整的请求链路日志记录
集成prometheus及pprof实现性能监控及性能问题排查
#### intialize
一个自动注入的配置注入及dao注入初始化，并暴露一个全局变量，记录模块信息
##### config（配置）
以nacos作为配置中心及本地配置文件（用于本地发开），支持toml格式的配置文件，
支持dev，test，prod环境本，启动命令区分
仅需配置nacos,后续配置均从nacos拉取及自动更新
```toml
Module = "user"
NoInit = ["Apollo","Etcd"]

[nacos]
Addr = "nacos.default"
Group = "hoper"
Watch  = false

[dev]
LocalConfigName = "local.toml"
[test]
NacosTenant = "b5b476ab-774f-4637-a8bf-e915635b4b24"
[prod]
NacosTenant = ""

```
##### dao（数据访问对象）
dao的注入基于配置，一个实现了`type Generate interface {Generate() interface{}}`的配置会自动生成一个dao
会根据dao结构体的字段名自动注入这个dao
#### context
```go
type Ctx struct {
	context.Context
	TraceID string
	*Authorization
	*DeviceInfo
	request.RequestAt
	Request *http.Request
	grpc.ServerTransportStream
	Internal string
	*log.Logger
}

```
一个通用的上下文，一个请求会生成一个context，贯穿整个请求，context记录原始请求上下文，请求时间，客户端信息，权限校验信息，及负责判断是否内部调用，
及附带唯一traceId的日志记录器
其中权限校验采用jwt，具体的校验模型采用接口，可供使用方自定义
```go
type AuthInfo interface {
	IdStr() string
}

type Authorization struct {
	AuthInfo     `json:"auth"`
	IdStr        string `json:"-" gorm:"-"`
	LastActiveAt int64  `json:"lat,omitempty"`
	ExpiredAt    int64  `json:"exp,omitempty"`
	LoginAt      int64  `json:"iat,omitempty"`
	Token        string `json:"-"`
}
```
一个典型的例子
```go
type AuthInfo struct {
	Id     uint64     `json:"id"`
	Name   string     `json:"name"`
	Role   Role       `json:"role"`
	Status UserStatus `json:"status"`
}
```

