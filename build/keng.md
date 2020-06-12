## 用其他主机docker login登录Harbor仓库报错
```bash
Error response from daemon: Get https://192.168.30.24/v2/: dial tcp 192.168.30.24:443: connect: connection refused
 vim /etc/docker/daemon.json
{
        "registry-mirrors": ["http://hoper.xyz"],
        "insecure-registries": ["192.168.xx.xx"]
}
restart docker
```
## Error loading config file XXX.dockerconfig.json - stat /home/XXX/.docker/config.json: permission denied
```
    这是因为docker的文件夹的权限问题导致的，处理办法如下，执行：
    
    sudo chown "$USER":"$USER" /home/"$USER"/.docker -R
    
    sudo chmod g+rwx "/home/$USER/.docker" -R
```

## Temporary failure in name resolution 错误
```bash
/etc/hosts
127.0.0.1       localhost.localdomain localhost
vim /etc/resolv.conf
nameserver   xxx
nameserver   xxx
```

## IDEA中总模块名与java中maven模块名冲突
改总模块名

## Java搞了半天缺依赖
```$xslt
pom中只有test，少了
 <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter</artifactId>
 </dependency>
```
## springboot管理普通类
@Component，@Autowired，@PostConstruct，init()

## java调go 远程主机强迫关闭了一个现有的连接。
建channel的时候少了usePlaintext()

## go调java 远程主机强迫关闭了一个现有的连接。
[https://github.com/grpc/grpc-java/issues/6011]
windows问题
So the problem is just the shutdown of the connection, which is not actually a problem.

## 调用wsl2上的grpc服务
监听地址应为0.0.0.0,不能是127.0.0.1

## redis MISCONF Redis is configured to save RDB snapshots, but it is currently not able to persist on disk. Commands that may modify the data set are disabled, because this instance is configured to report errors during writes if RDB snapshotting fails (stop-writ
强制把redis快照关闭了导致不能持久化
1.Redis客户端执行：config set stop-writes-on-bgsave-error no
2.修改redis.conf文件，stop-writes-on-bgsave-error=yes修改为stop-writes-on-bgsave-error=no

## Unsupported class file major version 57
升级到最新gradle

## Idea SpringBoot工程提示 "Error running 'xxxx'": Command line is too long.
1、找到workspace.xml文件

2、在<component name="PropertiesComponent">中添加<property name="dynamic.classpath" value="true" />一行

## spring.cloud.nacos.config.server-addr不生效
新建bootstrap.properties文件,该配置必须在启动加载配置文件中

## Gradle kotlin Springboot多模块导致无法引用kotlin的类文件(BootJar)
BUG项目 由于以Kotlin和Springboot中的多模块内容进行编写架构中，
发现 bootJar我用kotlin编写的jar包无法被正常的引用到，通过Gradle和SpringBoot项目下的Issue询问 ，
发现是由于Springboot插件，由于我的子模块集成了父容器的SpringBoot插件，导致 默认关闭了jar任务。原因连接[https://docs.spring.io/spring-boot/docs/2.1.4.RELEASE/gradle-plugin/reference/html/#managing-dependencies-using-in-isolation]
在你的子模块内容开发jar包任务如下
如果是Grovvy管理的：
```groovy
jar {
	enabled = true
}
```


如果是kotlin的kts管理的：
```kotlin
tasks.getByName<Jar>("jar") {
	enabled = true
}

```
[https://github.com/spring-projects/spring-boot/issues/16689]
[https://github.com/gradle/gradle/issues/9310]

## idea go debug 枚举值不显示值
右键 as Hex as Decimal as Binaty

## no Go source files
手动添加go.mod文件google.golang.org/protobuf（不知道有没有效）
idea 文件直接import需要的包，然后sync packages of 

##  no Go files in
main包路径不对
go install github.com/golang/protobuf
can't load package: package github.com/golang/protobuf: no Go files in E:\gopath\src\github.com\golang\protobuf
go install github.com/golang/protobuf/protoc-gen-go

## nacos post请求报参数错误
手动加请求头Content-Type:application/x-www-form-urlencoded

## 编译postwoman报错
清除npm缓存，npm i
好吧，postwoman那界面我受不了，内存大点就大点吧，其实apipost是真好用

## go Type.NumIn不一致
value := reflect.ValueOf(func)
value.Type().Method(j).Type.NumIn() 3 //方法第一个参数为接收器
value.Method(j).Type().NumIn() 2

## java11 Error: -p requires module path specification
在启用module的情况下，idea启动shorten command line 选user-local default: @argfile 会报这个错
选none不报错
原因是短命令行有个-p用来指定模块路径，然而并没有设置
If your command line is too long (OS is unable to create a process from command line of such length), IDEA provides you means to shorten the command line. The balloon should contain a link where you can switch between `none`|`classpath file`|`JAR-Manifest`|`args file (only available for jdk 9)`. If you run with any option but `none`, the command line would be shortened and OS should be able to start the process.

Do you still have a balloon when use one of the suggested (except `none`, please) options? If so, please enable debug option #com.intellij.execution.runners.ExecutionUtil (Help | Debug Log Settings), repeat running tests and attach idea.log (Help | Show log)

# 'Java SE 11' using tool chain : 'JDK 8 (1.8)'
 sourceCompatibility = JavaVersion.VERSION_11
 
# gradle java moudle
设置moduleName
 inputs.property("moduleName", moduleName)
  options.compilerArgs = listOf(
    "--module-path", classpath.asPath)
  classpath = files()
其他待解决问题：
vertx的依赖问题
同时读取io.netty
slf4j.log4j12 的依赖问题
错误: 模块 jvm 同时从 slf4j.log4j12 和 log4j 读取程序包 org.apache.log4j

# Every derived table must have its own alias
在做多表查询，或者查询的时候产生新的表的时候会出现这个错误：Every derived table must have its own alias（每一个派生出来的表都必须有一个自己的别名）。

# windows OpenSSH WARNING: UNPROTECTED PRIVATE KEY FILE!
ssh-keygen -t rsa

$env:username
更改文件所有者

vim /etc/ssh/sshd_config
AuthorizedKeysFile   .ssh/authorized_keys   //公钥公钥认证文件
RSAAuthentication yes
PubkeyAuthentication yes   //可以使用公钥登录

vim ~/.ssh/authorized_keys

service sshd restart

# nginx nginx: [emerg] unexpected "}" in
空格与制表符，nginx每行配置不支持空格开头