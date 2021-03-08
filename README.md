# hoper

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
                                             	pick.Api(func() interface{} {
                                             		return pick.Method(http.MethodPost).
                                             			Title("用户注册").
                                             			Version(2).
                                             			CreateLog("1.0.0", "jyb", "2019/12/16", "创建").
                                             			ChangeLog("1.0.1", "jyb", "2019/12/16", "修改测试")
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
https://github.com/kubernetes/ingress-nginx 已经用openresty了，哈哈哈
