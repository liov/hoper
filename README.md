# hoper

# alg_lang
各个语言的语法学习及刷LeetCode的解答

# build
开发及部署需要的操作总结脚本汇总，开发过程中遇到的坑总结
## config
    - nginx的配置
    - tls证书
    - rust cargo的配置
    - gradle 配置
## env
    开发所需环境的安装步骤，及安装过程中的问题解决方法
## k8s
    各种k8s yaml配置
## shell
    开发和部署用的到的脚本，docker和kubectl的命令

# client
## desktop
  flutter desktop尝鲜
## flutter
    hoper的移动端，flutter开发，Getx状态管理，组件化开发，开发了闪屏页，登录注册功能，瞬间列表，瞬间详情，发布照片瞬间，点赞评论等，grpc调用服务端接口，尝试了简单的dart ffi，调用go，rust交叉编译的动态库，webview嵌入，启动本地服务存储加载页面实现页面热更新，集成sqlite,hive键值对存储
## h5
    hoper的h5版本，vue3+typescript，开发了登录注册功能，瞬间列表，瞬间详情，发布瞬间，点赞收藏评论等

# server
## go
### lib
#### protobuf
通用的protobuf定义，通用的请求和返回定义，用到的工具protobuf定义
- generate 配合tools工具及第三方工具，生成自定义的http，grpc，graphql服务及api文档，
#### tiga框架
    - intialize 配置初始化及dao对象自动注入，
    -- 自由组合的conf和dao结构体
    -- 可选的配置中心apollo，etcd，nacos及本地配置,默认nacos
    -- 可选的dao对象，DB,pebble,redis,本地cache,mail,nsq,etcd.Client,kafka,elastic,badger
    本地的基本配置写入nacos配置中心地址，配置默认toml格式，main函数第一句`defer initialize.Start(uconf.Conf, udao.Dao)()`即可拉取nacos的配置自动注入配置，并根据配置注入创建dao对象
    - pick
        
    - server
        通用的对外暴露服务的框架，编写一套业务代码对外暴露http，grpc，graphql接口，集成trace记录，链路追踪，prometheus指标监控，pprof，debug，优雅关闭
#### tools
    protoc plugin，配合自定义protobuf文件，生成服务所需的源文件
    - protoc-gen-enum 针对枚举，生成中文的String();生成graphql所需的MarshalGQL，UnmarshalGQL，针对错误枚举生成error，GRPCStatus方法
    - protoc-gen-go-patch，fork 自protoc-gen-go-patch，protoc的插件，通过分析protoc生成的源码转为ast，修改ast输出源码来满足自定义属性生成需要的代码，主要用于结构体自定义标签
    - protoc-gen-grpc-gin，参照protoc-gen-grpc-gateway，生成基于gin的路由服务，用于grpc的业务函数对外外暴露http服务
#### 
    - context 一个通用的带有登录验证，trace记录，请求元信息的上下文定义，登录验证采用jwt，记录最基本的登录信息，
    
    


## 大杂烩式的项目
```
    -awesome
            -algorithms_languages  总结的各个语言的语法，坑和黑魔法，当然关于go的最多
            -os_compilers_graphics 操作系统，命令相关
            -protocols_softs_tools 协议，常用的软件和工具的安装操作
            -reproduceds_blogs 文章
   -build 本项目所用到的软件安装，发布脚本，基本配置，以及在项目过程中遇到的坑
   -client
            -flutter app客户端，采用flutter，添加原生，lua热更ui及热更业务逻辑，一种动态开发跨端app的模板
            -tarojs 小程序客户端， TODO 
            -vhoper h5版web
   -proto 本项目用到的proto文件
   -server  跨语言grpc模块
            -go go语言的模块 grpc+restful+graphql,每个模块都能拆出来当个库
                -intialize 初始化模块，负责conf和dao的自动注入，
                           支持开发环境的本地配置及测试生产环境配置中心配置拉取及监测，及配置中心的服务注册，
                            自动注入如DB,Redis，Nsq等客户端，
                -tools 项目需要用到的工具
                        -protoc-gen-enum 针对中文环境，proto文件enum枚举的补充生成，分错误枚举和补充枚举，
                                            分别生成Error()和Sting(),及其他方法
                        -protoc-gen-go-patch 对开源项目的pr，虽然作者没通过又自己写了一份，也算贡献了个issue
                                             通过对protoc-gen-go-grpc生成的源码解析为token再加工生成具有自定义tag的
                                             结构体及一些其他操作
                        -protoc-gen-grpc-gin 仿照protoc-gen-grpc-gateway，gateway太重了
                                             重写了基于gin的grpc->restful转接，当然仅限本地调用，并不算gateway
                -utils  工具包，有些可能已经达到库的级别
                       -log 对zap的包装，扩展输出方式，实现grpclog，另外对gorm log的实现放在dao.db.gorm包下
                            自定义输出格式，可配合kibana做日志查询
                       -net.http
                                -client 一个http请求客户端，定义请求与返回结构体即可方便的发送请求,包含返回注入及错误处理
                                         及日志记录
                                        ```go
                                            req := &OperatorReq{
                                            		PageNo:       1,
                                            		PageNum:      10,
                                            		PlatformType: 1,
                                            	}
                                            	res:=&OperatorPublicList{}
                                            	err := NewRequest(`http://hoper.xyz/api/list`,http.MethodPost,req).
                                            		SetHeader("Auth","e30=").HTTPRequest(CommonResponse(res))
                                            	if err!=nil{
                                            		t.Fatal(err)
                                            	}
                                        ```
                                -api swagger文档生成
                                -pick 一个基于自动注入的api框架
                                        -tree,router httprouter的基础上，不同请求方法的Handler采用切片存储，节省空间，提高性能
                                        -service 方法自动注入，分析方法名，入参出参自动注册路由，支持pick，gin，fiber，grpc_service
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
                                        -api_swagger swagger文档生成及访问，配合http.api
                                        -api_md markdown文档生成
                                        -grpc_gateway 针对grpc service 的路由注册
                                -gin 对gin框架的补充和扩展，http.Handler->gin.Handler的转换，上传导出的通用接口。oauth2的实现
                                -requset.binding 对入参与结构体的绑定注入
                                -tailmon 仍未想好命名的框架，与intialize搭配使用，甚至应该与intialize一个目录
                                          负责服务的启动，gin，grpc的整合，同时支持grpc调用和restful调用，
                                            整合prometheus，opencensus，集成性能监控和链路追踪，同时利用Context传递值，
                                            在进去业务方法前替换request中的Context，在业务方法中从context提取必要信息，
                                            同时整合log及opencensus，日志携带traceID，可跟踪查询每次请求的完整调用链日志
            -kotlin kotlin语言模块 springboot+grpc，quarkus的尝试
            -lua 结合openresty，后端的大前端，限流，黑白名单
            -rust rust的grpc模块
   -serverless serverless尝试
   -tools 平时写的小工具
```
## 架构设想
istio+k8s
