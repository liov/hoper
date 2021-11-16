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
测试不是github.com/xuri/excelize/v2的锅

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

# Property '$toast' does not exist on type 'Login'
ts问题 识别不了类型
@vue/cli-plugin-typescript 版本回退

# vue3.0 响应式Map在模板中获取不到值
``` vue
<template>
  <div>
    <van-action-sheet
      :show="show.moreShow"
      :actions="report.actions"
      cancel-text="取消"
      close-on-click-action
      @click="show.moreShow = !show.moreShow"
      teleport="#app"
    >
    </van-action-sheet>
    <van-dialog
      :show="report.show"
      title="举报"
      show-cancel-button
      @cancel="report.show = !report.show"
      teleport="#app"
    >
      <van-field name="radio">
        <template #input>
          <van-radio-group
            v-model="report.checked"
            direction="horizontal"
            @change="remark"
          >
            <van-radio name="1" shape="square">色情暴力</van-radio>
            <van-radio name="2" shape="square">侮辱谩骂</van-radio>
            <van-radio name="3" shape="square">政治政策</van-radio>
            <van-radio name="255" shape="square">其他原因</van-radio>
          </van-radio-group>
        </template>
      </van-field>
      <van-field
        v-if="report.field"
        v-model="report.message"
        rows="1"
        autosize
        label="备注"
        type="textarea"
        placeholder="请输入举报内容"
      />
    </van-dialog>
    <van-pull-refresh
      v-model="pullDown.refreshing"
      :success-text="pullDown.successText"
      @refresh="onRefresh"
    >
      <van-list
        :loading="state.loading"
        :finished="state.finished"
        finished-text="没有更多了"
        @load="onLoad"
      >
        <van-cell v-for="(item, index) in state.list">
          <template #default>
            <van-skeleton title avatar round :row="3" :loading="state.loading">
              <div class="moment" v-if="show.listShow">
                <div class="auth">
                  <img
                    class="avatar"
                    :src="state.map.get(item.userId).avatarUrl"
                  />
                  <span class="name">{{
                    state.map.get(item.userId).name
                  }}</span>
                  <span class="time">{{ $date2s(item.createdAt) }}</span>
                </div>
                <div class="content">
                  <van-field
                    v-model="item.content"
                    rows="1"
                    :autosize="{ maxHeight: 200 }"
                    readonly
                    type="textarea"
                  >
                    <template #extra>
                      <div class="arrow">
                        <van-icon name="arrow-down" />
                      </div>
                    </template>
                  </van-field>
                </div>
                <lazy-component class="imgs" v-if:="item.images">
                  <van-image
                    width="100"
                    height="100"
                    v-for="(img, idx) in item.images.split(',')"
                    :src="img"
                    lazy-load
                    class="img"
                    @click="preview(idx, item.images)"
                  />
                </lazy-component>
              </div>

              <van-row>
                <van-col
                  span="6"
                  class="action"
                  @click="show.moreShow = !show.moreShow"
                  ><van-icon name="more-o"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon
                    :name="item.collect ? 'star' : 'star-o'"
                    :color="item.collect ? '#F6DF02' : ''"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon name="comment-o"
                /></van-col>
                <van-col span="6" class="action"
                  ><van-icon
                    :name="item.likeId > 0 ? 'like' : 'like-o'"
                    :color="item.likeId > 0 ? '#D91E46' : ''"
                    @click="like(index)"
                /></van-col>
              </van-row>
            </van-skeleton>
          </template>
        </van-cell>
      </van-list>
    </van-pull-refresh>
  </div>
</template>

<script lang="ts" setup>
import axios from "axios";
import { ObjMap } from "@/plugin/utils/user";
import { ImagePreview } from "vant";
import { reactive, ref } from "vue";

const pageNo = ref(1);
const pageSize = 10;
const userM = new ObjMap();
const state = reactive({
  loading: false,
  finished: false,
  list: Array.from(new Array(pageSize), (v, i) => {
    return { id: i };
  }),
  map: new Map(),
});

const pullDown = reactive({
  successText: "刷新成功",
  refreshing: false,
});
const show = reactive({
  listShow: false,
  moreShow: false,
  shareShow: false,
});

const report = reactive({
  show: false,
  actions: [
    { name: "不喜欢" },
    {
      name: "举报",
      callback: () => (report.show = !report.show),
    },
    {
      name: "删除",
      color: "#D91E46",
    },
  ],
  checked: false,
  field: false,
  message: "",
});
//mounted() {}
const onLoad = async () => {
  state.loading = false;
  // 异步更新数据
  const res = await axios.get(
    `/api/v1/moment?pageNo=${pageNo.value}&pageSize=${pageSize}`
  );
  if (res.data.code !== 0) {
    this.$toast.fail(res.data.message);
    state.finished = true;
  }
  const data = res.data.details;
  if (state.pageNo == 1) {
    state.list = data.list;
  } else {
    state.list = state.list.concat(data.list);
  }
  userM.appendMap(data.users);
  for (let user of data.users) {
    state.map.set(user.id, user);
  }
  state.loading = false;
  show.listShow = true;
  pageNo.value++;
  if (data.list.length < pageSize) state.finished = true;
};
const preview = (idx: number, images: string) => {
  ImagePreview({
    images: images.split(","),
    startPosition: idx,
    closeable: true,
  });
};
const onRefresh = () => {
  pullDown.refreshing = true;
  pageNo.value = 1;
  onLoad().catch(() => {
    pullDown.successText = "刷新失败";
  });
  pullDown.refreshing = false;
};
const remark = (name: string) => {
  console.log(name);
  if (name === "255") {
    report.field = true;
  }
};
const like = async (idx: number) => {
  console.log(state.list[idx]);
  const api = `/api/v1/action/like`;
  const id = state.list[idx].id;
  const likeId = state.list[idx].likeId;
  if (likeId > 0) {
    await axios.delete(api, { data: { id: likeId } });
    state.list[idx].likeId = 0;
  } else {
    state.list[idx].likeId = await axios.post(api, {
      refId: id,
      type: 1,
      action: 2,
    });
  }
};
</script>

<style scoped lang="less">
.moment {
  @20px: 20px;
  @avatar: 30px;
  .name {
    left: 60px;
    position: absolute;
  }

  .time {
    position: absolute;
    right: @20px;
  }
  .content {
    width: 100%;
    h3 {
      margin: 0;
      font-size: 18px;
      line-height: 20px;
    }

    .arrow {
      position: absolute;
      bottom: 16px;
      right: 0;
    }

    .van-multi-ellipsis--l3 {
      margin: 13px 0 0;
      font-size: 14px;
      line-height: 20px;
    }
  }

  .avatar {
    flex-shrink: 0;
    width: @avatar;
    height: @avatar;
    border-radius: 40px;
    position: relative;
    margin: 0 16px;
  }
  .imgs {
    padding: 0 11px;
  }
  .img {
    margin: 5px 5px;
  }
  .action {
    text-align: center;
  }
}
</style>
```
# The import path must contain at least one forward slash ('/') character.
See https://developers.google.com/protocol-buffers/docs/reference/go-generated#package for more information.
--go_out: protoc-gen-go: Plugin failed with status code 1.
Before
option go_package = ".;openapiconfig";
After
option go_package = "github.com/liov/hoper/server/go/lib/protobuf/utils/proto/openapiconfig";

# node_modules\deasync: Command failed
Exit code: 1
Command: node ./build.js
Arguments:
Directory: .\node_modules\deasync
Output:
node:events:346
      throw er; // Unhandled 'error' event
      ^
Error: spawn node-gyp.cmd ENOENT
    at Process.ChildProcess._handle.onexit (node:internal/child_process:282:19)
    at onErrorNT (node:internal/child_process:480:16)
    at processTicksAndRejections (node:internal/process/task_queues:81:21)
Emitted 'error' event on ChildProcess instance at:
    at Process.ChildProcess._handle.onexit (node:internal/child_process:288:12)
    at onErrorNT (node:internal/child_process:480:16)
    at processTicksAndRejections (node:internal/process/task_queues:81:21) {
  errno: -4058,
  code: 'ENOENT',
  syscall: 'spawn node-gyp.cmd',

npm install -g node-gyp

# emmit事件注册无效

vue下一次渲染的created函数先于上一个组件的unmounted执行，应该在mounted中注册事件

# wsl2和windows的网络问题
监听127.0.0.1和0.0.0.0时，natstat显示的Local Address是不同的，
    localhost:50051 [::]:50051
0.0.0.0 不能ping通，代表本机所有的IP地址；
    监听127.0.0.1，创建Socket，那么用本机地址建立tcp连接不成功，反过来也是如此；也就是，监听时采用的地址为192.168.0.1，就只能用192.168.0.1进行连接。
    而监听0.0.0.0创建Socket，那么无论使用127.0.0.1或本机ip都可以建立tcp连接,也就是不论通过127.0.0.1或192.168.0.1、192.168.1.1都能连接成功。
    0.0.0.0建立tcp连接的时候也可以通过绑定IP_ADDR_ANY来实现。
IPv4 的环回地址是保留地址之一 127.0.0.1。尽管只使用 127.0.0.1 这一个地址，但地址 127.0.0.0 到 127.255.255.255 均予以保留。此地址块中的任何地址都将环回到本地主机中。此地址块中的任何地址都绝不会出现在任何网络中。
首先我们来讲讲127.0.0.1，172.0.0.1是回送地址，localhost是本地DNS解析的127.0.0.1的域名，在hosts文件里可以看到。

一般我们通过ping 127.0.0.1来测试本地网络是否正常。其实从127.0.0.1~127.255.255.255，这整个都是回环地址。这边还要

注意的一点就是localhost在了IPV4的是指127.0.0.1而IPV6是指::1。当我们在服务器搭建了一个web服务器的时候如果我们

监听的端口时127.0.0.1：端口号 的 时候，那么这个web服务器只可以在服务器本地访问了，在别的地方进行访问是不行的。

（127.0.0.1只可以在本地ping自己的，那么你监听这个就只可以在本地访问了）

  然后我们来讲讲0.0.0.0，如果我们直接ping 0.0.0.0是不行的，他在IPV4中表示的是无效的目标地址，但是在服务器端它表示

本机上的所有IPV4地址，如果一个服务有多个IP地址（192.168.1.2和10.1.1.12），那么我们如果设置的监听地址是0.0.0.0那

么我们无论是通过IP192.168.1.2还是10.1.1.12都是可以访问该服务的。在路由中，0.0.0.0表示的是默认路由，即当路由表中

没有找到完全匹配的路由的时候所对应的路由。

# Error: unknown file extension .ts
移除"type": "module" in package.json 会导致import报错 正确做法，  
1.移除"type": "module" in package.json  
2.tsconfig "module": "commonjs"

# typescript

## Locally in your project.
npm install -D typescript
npm install -D ts-node

## Or globally with TypeScript.
npm install -g typescript
npm install -g ts-node

## Depending on configuration, you may also need these
npm install -D tslib @types/node

# TS1378: Top-level 'await' expressions are only allowed when the 'module' option is set to 'esnext' or 'system', and the 'target' option is set to 'es2017' or higher.
"type": "module" in package.json
yarn add --dev ts-node yarn add --dev
typescript yarn add --dev @types/node
node --loader=ts-node/esm

# Error: ERR_MODULE_NOT_FOUND D:\hoper\tools\serverless\src\utils\random D:\hoper\tools\serverless\src\ts\db_mock.ts module
- import {randomNum, lottery} from '../utils/random';
+ import {randomNum, lottery} from '../utils/random.js';

# flutter PageView + TabView 不能连续滑动
TabView是在PageView的基础上封装的
其实可以通过自定义PageViewTabClampingScrollPhysics,根据边界，滑动量偏移量来判断滑动，根据PageView的Controller.animateToPage跳转到
相应页面，处理得当体验是相当好的，但是有个小bug，PageView滑动太快会直接跳过TabView，还有性能差
Listener({
  Key key,
  this.onPointerDown, //手指按下回调
  this.onPointerMove, //手指移动回调
  this.onPointerUp,//手指抬起回调
  this.onPointerCancel,//触摸事件取消回调
  this.behavior = HitTestBehavior.deferToChild, //在命中测试期间如何表现
  Widget child
})或者GestureDetector似乎也能实现
所以就只能牺牲体验，把PageView设置为不可滑动，由底部导航栏导航

# flutter 切换选项卡会rebuild
StatelessWidget 目前无解
StatefulWidget AutomaticKeepAliveClientMixin
所以Get做全局的状态管理就好了,局部的并不需要

# 没有root账号但有root权限的机器上Minikube
sudo minikube
sudo kubectl
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add apisix https://charts.apiseven.com
helm repo update
helm pull apisix/apisix
sudo helm install apisix ./apisix-0.3.5.tgz \
  --set admin.allow.ipList="{0.0.0.0/0}" \
  --namespace ingress-apisix
sudo helm install apisix-ingress-controller ./apisix-ingress-controller \
  --set image.tag=dev \
  --set config.apisix.baseURL=http://apisix-admin:9180/apisix/admin \
  --set config.apisix.adminKey=edd1c9f034335f136f87ad84b625c8f1 \
  --namespace ingress-apisix

# spec.ports[0].nodePort: Invalid value: 80: provided port is not in the valid range. The range of valid ports is 30000-32767
minikube native
# 无效
minikube start --extra-config=apiserver.GenericServerRunOptions.ServiceNodePortRange=1-10000
# 有效
sudo vim /etc/kubernetes/manifests/kube-apiserver.yaml
command 下添加 --service-node-port-range=1-65535 参数
kill 掉 kube-apiserver

# 编译在docker alpine linux中可用的go程序
CGO_ENABLED=0 go build
----------------------
go build -tags netgo
-------------------
```Dockerfile
FROM docker.io/golang:alpine

RUN echo "https://mirror.tuna.tsinghua.edu.cn/alpine/v3.14/main" > /etc/apk/repositories

RUN apk add --no-cache gcc musl-dev

```
docker build -t go-build:1.0 .
docker run -e "GOPROXY=https://goproxy.io" -it --rm -v `pwd`:/root/src -w /root/src  go-build:1.0  go build github.com/Kong/go-pluginserver

/usr/local/go/pkg/tool/linux_amd64/link: running gcc failed: exec: "gcc": executable file not found in $PATH
------------------------------------------------------------------------------------------------------------------------------------------------------
mkdir /lib64
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 \
&& ln -s /usr/lib/libGraphicsMagickWand.so.2.9.4 /lib/libGraphicsMagickWand-Q16.so.2 \
&& ln -s /usr/lib/libGraphicsMagick.so.3.21.0 /lib/libGraphicsMagick-Q16.so.3

# go.info.runtime.firstmoduledata: relocation target go.info.github.com/liov/hoper/server/go/lib/utils/reflect.moduledata not defined
https://github.com/golang/go/issues/46777
这是个从不被有意支持的功能，不建议使用
//go 1.17报错，要带完整包名让编译器找的到
////go:linkname Firstmoduledata runtime.firstmoduledata
//go:linkname reflecti.Firstmoduledata runtime.firstmoduledata

# win11 emulator 打不开
HAXM 和Hyper-v冲突  
无法卸载HAXM,windows可选功能关闭Hyper-v无效
管理员身份
bcdedit /set hypervisorlaunchtype off 重启
bcdedit /set hypervisorlaunchtype auto

没试过android-sdk卸载HAXM打开Hyper-v是否可行

# flutter pub run . Could not find a file named "pubspec.yaml" in
flutter pub get
flutter pub run

# flutter pub run flutter_native_splash:create android12无效
天坑，不看源码还不知道，flutter_native_splash是根据build.gradle判断编译SKD版本的，判断方法简单粗暴截取转整型，后面有注释识别不了
去掉注释