# hoper

# algo_lang
各个语言的语法学习,总结的坑和黑魔法,刷LeetCode的解答
## lua
 专注于openresty开发和apisix插件
## python
一些gui samples
# build
开发及部署需要的操作总结脚本汇总，开发过程中遇到的坑总结
## config
- nginx的配置
- rust cargo的配置
- gradle 配置
## env
### dev
开发所需工具的安装步骤，及安装过程中的问题解决方法，及使用方法
### shell
开发所需脚本
## k8s
各种工具部署的k8s 配置
### app
个人项目配置drone所需的Dockerfile，k8s配置文件模板、脚本
### tools
k8s环境中所需工具的安装方法，helm & yaml config,及使用方法
### tpl
k8s配置模板

# client
## desktop
flutter desktop尝鲜
## flutter
 hoper的移动端
flutter开发，Getx状态管理，组件化开发
开发了闪屏页，登录注册功能，瞬间列表，瞬间详情，发布照片瞬间，点赞评论等
grpc调用服务端接口
尝试了简单的dart ffi，调用go，rust交叉编译的动态库，webview嵌入，启动本地服务存储加载页面实现页面热更新
集成sqlite,hive键值对存储
## h5
hoper的h5版本，vue3+typescript，开发了登录注册功能，瞬间列表，瞬间详情，发布瞬间，点赞收藏评论等
## mini
taro+vue 开发的hoper小程序，编译的h5版本将取代h5版本

# proto
hoper项目客户端，服务端通用的proto定义，用来生成go，rust，java，dart，js grpc源文件

# server
## go
### lib
 一个通用且强大的库
#### context
一个通用的带有登录验证，trace记录，请求元信息的上下文定义，登录验证采用jwt，记录最基本的登录信息
#### intialize
- intialize 配置初始化及dao对象自动注入，
- 自由组合的conf和dao结构体
- 可选的配置中心apollo，etcd，http文件，nacos及本地配置,默认nacos
- 可选的dao对象，DB,pebble,redis,本地cache,mail,nsq,etcd.Client,kafka,elastic,badger

本地的基本配置写入nacos配置中心地址，配置默认toml格式
main函数第一句`defer initialize.Start(uconf.Conf, udao.Dao)()`即可拉取nacos的配置自动注入配置，并根据配置注入创建dao对象
#### pick
 一个基于自动注入的api框架，可自动注册路由，自动注入入参，可兼容gin fasthttp
    
- tree,router httprouter的基础上，不同请求方法的Handler采用切片存储，节省空间，提高性能
- service 方法自动注入，分析方法名，入参出参自动注册路由，支持pick，gin，fiber，grpc_service
```go
func (*UserService) Add(ctx context.Context, req *model.SignupReq) (*response.TinyRep, error) {
//对于一个性能强迫症来说，我宁愿它不优雅一些也不能接受每次都调用
pick.Api(func() {
pick.Method(http.MethodPost).
Title("用户注册").
Version(2).
CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试").End()
})

return &response.TinyRep{Message: "测试"}, nil
}
```
- api_swagger swagger文档生成及访问，配合http.api
- api_md markdown文档生成
- grpc_gateway 针对grpc service 的路由注册
#### protobuf
通用的protobuf定义，通用的请求和返回定义，用到的工具protobuf定义
 - generate 配合tools工具及第三方工具，生成自定义的http，grpc，graphql服务及api文档，
#### tiga
通用的对外暴露服务的框架，编写一套业务代码对外暴露http，grpc，graphql接口，集成trace记录，链路追踪，prometheus指标监控，pprof，debug，优雅关闭
- 与intialize，context搭配使用
- 负责服务的启动，gin，grpc的整合，同时支持grpc调用和restful调用，
- 整合prometheus，opencensus，集成性能监控和链路追踪，同时利用Context传递值，
- 在进去业务方法前替换request中的Context，在业务方法中从context提取必要信息，
- 同时整合log及opencensus，日志携带traceID，可跟踪查询每次请求的完整调用链日志
#### tools
protoc plugin，配合自定义protobuf文件，生成服务所需的源文件
- protoc-gen-enum 针对枚举，生成中文的String();生成graphql所需的MarshalGQL，UnmarshalGQL，针对错误枚举生成error，GRPCStatus方法
- protoc-gen-go-patch，fork 自protoc-gen-go-patch，protoc的插件，通过分析protoc生成的源码转为ast，修改ast输出源码来满足自定义属性生成需要的代码，主要用于结构体自定义标签
- protoc-gen-go-graphql，用于生成graphql的定义文件
- protoc-gen-grpc-gin，参照protoc-gen-grpc-gateway，生成基于gin的路由服务，用于grpc的业务函数对外外暴露http服务
#### utils
各种工具类
## kotlin
### protobuf
protoc生成的源码单独成一个模块
### quarkus
quarkus框架
### user
- spring+vertx
- spring+grpc
## rust
 基于tonic的grpc
# tools
平时写的小工具
## python
爬虫收集，人工智能应用
## script
ts&js，自动化脚本
## server
爬虫，处理工具
## 架构设想
istio+k8s
