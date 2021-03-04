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

# nacos k8s 部署503
[执行sql](https://github.com/alibaba/nacos/blob/b9ff53b49cec5ca7cf37736ebc9c1c2bb4a108a8/config/src/main/resources/META-INF/nacos-db.sql)

# docker修改/etc/docker/daemon.json后无法重启

vim /usr/lib/systemd/system/docker.service 删除冲突配置
systemctl daemon-reload
systemctl restart docker.service

# k8s.gcr.io
docker pull registry.aliyuncs.com/google_containers/<imagename>:<version>
docker tag registry.aliyuncs.com/google_containers/<imagename>:<version> k8s.gcr.io/<imagename>:<version>
```bash
eval $(echo ${images}|
        sed 's/k8s\.gcr\.io/anjia0532\/google-containers/g;s/gcr\.io/anjia0532/g;s/\//\./g;s/ /\n/g;s/anjia0532\./anjia0532\//g' |
        uniq |
        awk '{print "docker pull "$1";"}'
       )
for img in $(docker images --format "{{.Repository}}:{{.Tag}}"| grep "anjia0532"); do
  n=$(echo ${img}| awk -F'[/.:]' '{printf "gcr.io/%s",$2}')
  image=$(echo ${img}| awk -F'[/.:]' '{printf "/%s",$3}')
  tag=$(echo ${img}| awk -F'[:]' '{printf ":%s",$2}')
  docker tag $img "${n}${image}${tag}"
  [[ ${n} == "gcr.io/google-containers" ]] && docker tag $img "k8s.gcr.io${image}${tag}"
done
```

# spring + vertx 浏览器NOT Found
```yaml
server:
  port: 8090
```
去掉这个配置,我们只用spring的依赖注入,springmvc或者springwebflux会自动读取占用端口开启服务

# InteIIiJ IDEA Gradle 编码 GBK 的不可映射字符
tasks.withType(JavaCompile) {
    options.encoding = "UTF-8"
}

# Android编译时报错：More than one file was found with OS independent path lib/armeabi-v7a/libluajapi.so

packagingOptions {
        // pickFirsts:当出现重复文件，会使用第一个匹配的文件打包进入apk
        pickFirst 'lib/armeabi-v7a/libluajapi.so'
    }
    
# Android Execution failed for JetifyTransform

compileOptions{
        sourceCompatibility JavaVersion.VERSION_1_8
        targetCompatibility JavaVersion.VERSION_1_8
    }

# IDEA安装插件后打不开，插件木录  
${Home}\AppData\Roaming\JetBrains\IntelliJIdea2020.1\plugins

# Android打包动态库
```groovy
android {
    sourceSets {
        main {
            jniLibs.srcDirs = ['src/main/jniLibs']
        }
    }
}
dependencies {
    implementation fileTree(dir: 'lib', include: ['*.so'])
}

```

# 服务器被挂马
top
ps -ef|grep xxx
ls -l /proc/pid
crontab -l
crontab -r
rm xxx

#windows 文件夹删不掉 该项目不在 请确认该项目的位置
```bat
DEL /F /A /Q \\?\%1
RD /S /Q \\?\%1
```
拖着要删除东西拉到bat文件上

# cmd中文乱码
chcp 65001

# etcd 共用
使用apisix，最初想与k8s集群共用etcd，但是minikube中无法实现,应该是minikube部署在docker中，docker重启IP变了，证书不认了

# minikube The connection to the server localhost:8443 was refused - did you specify the right host or port? waiting for app.kubernetes.io/name=ingress-nginx pods: timed out waiting for the condition]
delete start
# pod内无法ping通svc
```bash
kubectl edit cm kube-proxy -n kube-system
mode:"ipvs"


cat >> /etc/sysctl.conf << EOF
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

kubectl  get pod -n kube-system | grep kube-proxy | awk '{print $1}' | xargs kubectl delete pod -n kube-system
```
# root用户读不到环境变量
sudo visudo

Defaults    !env_reset

# minikube diver=none minikube kubectl 无法使用
sudo /usr/local/bin/minikube 

# nodePort 80
vim /etc/kubernetes/manifests/kube-apiserver.yaml
command 下添加 --service-node-port-range=1-65535 参数

# js正则匹配失败
文件换行符CRLF -> LF

# go交叉编译的bug cannot find module for path 
正常编译可以，交叉编译就报包找不到(cannot find module for path github.com/360EntSecGroup-Skylar/excelize)
main里下划线导入不报找不到包(https://juejin.im/post/5d776830f265da03e05b3c45),内部包找不到了
cgo的锅
set CGO_ENABLED=1
测试不是github.com/360EntSecGroup-Skylar/excelize/v2的锅

应该是cgo的原因，但是那个项目里的包都是常见的包啊，难以定位哪里用了带cgo的包

交叉编译时，CGO_ENABLED=0是会自动忽略带cgo的包，这个有bug，1.14会修复[https://github.com/golang/go/issues/35873]

main包匿名导入提示找不到路径的包又不报这个错，报内部包的函数undefine
无法复现

排查了半天，真的让人哭笑不得
真的跟cgo有关
那个引用找不到路径的包的包多了个import "C"，不知道什么时候加上去的

---p1.go
package p1

import "C"
import github.com/user/p2

---go.mod
github.com/user/p2

# x86_64-w64-mingw32/bin/ld.exe: Error: export ordinal too large

Go tool argument -buildmode=exe

# Parameter 'xxx' implicitly has an 'any' type.
tsconfig.json添加"noImplicitAny": false，

或者 "strict": true,改为false

# postgres默认时间
时区调成上海后，设定默认时间'0001-01-01 00:00:00+08'总会自动变成'0001-01-01 00:00:00+08:05:43'::timestamp with time zone
加的时间不正常，试了几次，分界时间是1900年，加时不对用时间过滤的时候会有问题，从01年以后都是正常加8