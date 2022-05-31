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
ssh-keygen -t rsa -b 4096 -C "autossh" -f autossh_id_rsa

$env:username
更改文件所有者

vim /etc/ssh/sshd_config
AuthorizedKeysFile   .ssh/authorized_keys   //公钥公钥认证文件
RSAAuthentication yes
PubkeyAuthentication yes   //可以使用公钥登录

cat autossh_id_rsa.pub >> ~/.ssh/authorized_keys

service sshd restart
报错 sshd -T

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
k8s.gcr.io/pause:3.1
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
-tag ping service clusterIP
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
option go_package = "github.com/actliboy/hoper/server/go/lib/protobuf/utils/proto/openapiconfig";

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
minikube start --extra-config=apiserver.service-node-port-range=1-10000
--extra-config=apiserver.service-node-port-range=1-65536
-- 最新试的有效
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
docker run -e "GOPROXY=https://goproxy.io" -it --rm -v `pwd`:/app -w /app  go-build:1.0  go build github.com/Kong/go-pluginserver

/usr/local/go/pkg/tool/linux_amd64/link: running gcc failed: exec: "gcc": executable file not found in $PATH
------------------------------------------------------------------------------------------------------------------------------------------------------
mkdir /lib64
ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 \
&& ln -s /usr/lib/libGraphicsMagickWand.so.2.9.4 /lib/libGraphicsMagickWand-Q16.so.2 \
&& ln -s /usr/lib/libGraphicsMagick.so.3.21.0 /lib/libGraphicsMagick-Q16.so.3

# go.info.runtime.firstmoduledata: relocation target go.info.github.com/actliboy/hoper/server/go/lib/utils/reflect.moduledata not defined
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

# flutter Scaffold bottomSheet 高度不够，溢出底部
最外层加SafeArea

# github有,goproxy.io拉不到
set GOPROXY=https://goproxy.cn

# 无法读取源文件或磁盘

chkdsk E:/R

# OCI runtime create failed: container_linux.go:380: starting container process caused: process_linux.go:385: applying cgroup configuration for process caused: mkdir /sys/fs/cgroup/memory/docker/68044e2dec19557f55a8a8fe6ea035d3ff30325165787f2449f9e2ea17152542: cannot allocate memory: unknown

https://zhuanlan.zhihu.com/p/106757502
CentOS 7.6 3.10.x 内核下创建 memory cgroup 失败原因与解决方案
3.10.x 内核的 memory cgroup 实现中 kmem accounting 相关是 alpha 特性，实现上由许多 BUG 。具体我们的场景，开启 kmem accouting 后有两个问题：

SLAB 泄露，具体可以参见 Try to Fix Two Linux Kernel Bugs While Testing TiDB Operator in K8s | TiDB
memory cgroup 泄露，删除后不会被完全回收。又因为内核对 memory cgroup 总数有 65535 总数限制，频繁创建删除开启了 kmem 的 cgroup ，会让内核无法再创建新的 memory cgroup ，出现上面的错误。

从 RedHat 的 Issue 中，据说问题在 7.8 的 kernel-3.10.0-1075.el7 最终修复了。等 CentOS 7.8 发布后，也是一个选择。

升级内核版本到7.6以上即可解决

如果没办法立刻升级，那么就先重启释放内存吧

https://www.cnblogs.com/xzkzzz/p/9627658.html

rpm --import https://www.elrepo.org/RPM-GPG-KEY-elrepo.org
rpm -Uvh http://www.elrepo.org/elrepo-release-7.0-3.el7.elrepo.noarch.rpm
yum --enablerepo=elrepo-kernel install kernel-ml

sudo awk -F\' '$1=="menuentry " {print i++ " : " $2}' /etc/grub2.cfg
grub2-set-default 0

vim /etc/default/grub
GRUB_DEFAULT=0
grub2-mkconfig -o /boot/grub2/grub.cfg
reboot

# Gradle Error:java.lang.OutOfMemoryError: Java heap space.

Gradle设置JVM内存的方法有以下几种方式：

命令行启动任务时增加设置堆内存的参数

修改用户目录下的gradle.properties文件

Windows下路径一般为：C:\Users\用户名.gradle，如果该目录没有这个文件可以新建一个。

修改项目目录下的grable.properties文件（没有可以新建一个）

如果使用的是IDE插件也可以在IDE里面进行设置

如果修改的是gradle.properties文件，加入以下配置：

org.gradle.jvmargs=-Xmx2048m -XX:MaxPermSize=512m
如果是在命令行或IDE里面调整JVM内存分配，只需要使用：

distributionUrl=https\://services.gradle.org/distributions/gradle-7.3-bin.zip
org.gradle.jvmargs=-Xmx4096m -XX:MaxPermSize=4096m -XX:+HeapDumpOnOutOfMemoryError
不行再大

# flutter Does not match the generator used previously: Visual Studio 16 2019
CMake Error: Error: generator : Visual Studio 17 2022
Does not match the generator used previously: Visual Studio 16 2019
Either remove the CMakeCache.txt file and CMakeFiles directory or choose a different binary directory.
Exception: Unable to generate build files

flutter clean

# Duplicate GlobalKey detected in widget tree. 
Flutter在initState()初始化方法时使用包含context的Widget导致报错问题
```dart
 WidgetsBinding.instance?.addPostFrameCallback((timeStamp) {
      globalService.log.finest("WidgetsBinding.instance");
      W = MediaQuery.of(context).size.width;
    });
```
# Error: "linker 'cc' not found" when cross compiling a rust project from windows to linux using cargo

It turns out you need to tell cargo to use the LLVM linker instead. You do this by creating a new directory called .cargo in your base directory, and then a new file called config.toml in this directory. Here you can add the lines:

[target.aarch64-linux-android]

ar = "Android\\Sdk\\ndk-bundle\\toolchains\\llvm\\prebuilt\\windows-x86_64\\bin\\aarch64-linux-android-ar.exe"

linker = "C:\\Users\\mayn\\AppData\\Local\\Android\\Sdk\\ndk-bundle\\toolchains\\llvm\\prebuilt\\windows-x86_64\\bin\\aarch64-linux-android28-clang.cmd"

[target.armv7-linux-androideabi]

ar = "Android\\Sdk\\ndk-bundle\\toolchains\\llvm\\prebuilt\\windows-x86_64\\bin\\arm-linux-androideabi-ar.exe"

linker = "Android\\Sdk\\ndk-bundle\\toolchains\\llvm\\prebuilt\\windows-x86_64\\bin\\armv7a-linux-androideabi28-clang.cmd

Then building with the command cargo build --target=x86_64-unknown-linux-musl should work!

# rust android ld: error: unable to find library -lgcc

https://github.com/mozilla/rust-android-gradle/issues/75

@ncalexan actually make_standalone_toolchain.py is never called. The problem is NDK 23.1.7779620 has stop distributing libgcc . I got around this by hacking linker-wrapper.py to statically linking libgcc.a from NDK 22.1.7171670

args.remove("-lgcc")
args.append("/home/sto/android/ndk/22.1.7171670/toolchains/arm-linux-androideabi-4.9/prebuilt/linux-x86_64/lib/gcc/arm-linux-androideabi/4.9.x/libgcc.a")
published this plugin locally with
gradle publishToMavenLocal

then rebuild my project with
gradle cargoBuild

then tested it with

sto@big:~/workspace/gossipy$ adb push ./app/src/main/rust/target/armv7-linux-androideabi/debug/hello-rust /data/local/tmp
./app/src/main/rust/target/armv7-linux-androideabi/debug/hello-rust: 1 file pushed, 0 skipped. 250.7 MB/s (3264536 bytes in 0.012s)
sto@big:~/workspace/gossipy$ adb shell /data/local/tmp/hello-rust
Hello, rust world!
is there a better way to do this? the plugin should have a way of including static libraries

there should also be an option to change the linker-wrapper.py to not link -lgcc and use static version instead

# serde cannot find derive macro Deserializ in this scope
在新版本的serde中使用derive需要开启对应的features
serde = {version="1.0.132",features=["derive"]}

# rust main.rs 和 lib.rs 同时存在，main不能用lib failed to resolve: use of undeclared crate or module `rust_lib` create use of undeclared crate or module `rust_lib`
cargo.toml crate-type "lib" 补上
```toml
name = "rust_lib"
crate-type = ["lib", "staticlib", "cdylib"]
```

# protoc卡住
tags标签写错

# postgres删表 is being access by other users

SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE datname='mydb' AND pid<>pg_backend_pid();

# postgres数据库迁移 navicat 备份导入字符串部分会报错

pg_dump -U postgres -p 5432 -d test -f /home/postgres/test12.sql
psql -d test -U postgres -f test12.sql


postgres进行迁移可以使用psql，也可以使用postgres自带工具pg_dump和pg_restore.

命令：

- 备份

pg_dump -h 13.xx.xx.76 -U postgres -n "public" "schema" -f ./schema_backup.gz -Z 9

-h host，备份目标数据库的ip地址

-U 用户名（输入命令后会要求输入密码，也可以使用-w输入密码）

-n 需要导出的schema名称

-f 导出存储的文件

-Z 进行压缩（一般导出文件会占用很大的存储空间，直接进行压缩）

- 恢复

gunzip schema_backup.gz ./ （对导出的压缩文件解压）

psql -U postgres -f ./schema_backup >>restore.log

参数意义与导出一样

坑与tips：

版本，pg_dump的版本要高于目标备份数据库的版本（比如目标数据库是10.3， pg_dump要使用10.3或者10.4）

-Z 是pg_dump提供的压缩参数，默认使用的是gzip的格式，目标文件导出后，可以使用gunzip解压（注意扩展名，有时习惯性命名为.dump 或者.zip，使用gunzip时会报错，要改为.gz）

也可以针对指定的表进行导出操作：

pg_dump -h localhost -U postgres -c -E UTF8 --inserts -t public.t_* > t_taste.sql

--inserts 导出的数据使用insert语句

-c 附带创建表命令

## 比较骚
1.操作位置：迁移数据库源（旧数据库主机）

找到PostgreSql 的data目录   关闭数据库进程

打包 tar -zcvf pgdatabak.tar.gz data/

------------------------------------------------------------------

2.通过winScp 或者 CRT 等工具拷贝到    迁移目标源（新主机--需安装postgresql）  同样的data目录 关闭数据库进程

解压  tar -zxvf pgdatabak.tar.gz -C /usr/local/postgres/

重新授权 执行命令  chown -R postgres.postgres data/

# error:14094438:SSL routines:ssl3_read_bytes:tlsv1 alert internal error
nginx 不能用这种写法
set $tls D:/code/hoper/build/config/tls;

ssl_certificate     $tls/cert.pem;
ssl_certificate_key  $tls/cert.key;
只能
ssl_certificate    D:/code/hoper/build/config/tls/cert.pem;
ssl_certificate_key D:/code/hoper/build/config/tls/cert.key;

# nginx *384 stat() "/home/ubuntu/deploy/hoper/dist/index.html" failed (13: Permission denied)
/etc/nginx/nginx.conf中的第一行改为user root;

# Got permission denied while trying to connect to the Docker daemon socket at
sudo gpasswd -a ${USER} docker
newgrp docker

# ERROR: cached plan must not change result type
postgres修改字段长度后出现，断开重连

# ERROR: duplicate key value violates unique constraint "diary_pkey" (SQLSTATE 23505)
单纯清除数据的方法
– 清除所有的记录（有外键关联的情况下）
TRUNCATE TABLE questions CASCADE;

– 清除所有的记录，并且索引号从0开始
TRUNCATE TABLE questions RESTART IDENTITY CASCADE;

# The system cannot find the path specified
 MkDirAll()

# go build -o timepill timepill/main.go undefined: Diary
go build -o timepill tools/timepill

# postgres 查询user表 ERROR:  column "id" does not exist

SELECT EXISTS(SELECT id FROM user WHERE user_id = 100192773 LIMIT 1);
SELECT EXISTS(SELECT id FROM "user" WHERE user_id = 100192773 LIMIT 1);

# docker宿主机访问docker容器服务失败

## 猜测原因
因为docker的虚拟ip网段是172.17.*。*与局域网的ip网段172.17冲突了，所以有两种方式：

解决方法：

一、

修改docker网卡信息，将网段改为与局域网不同的即可
```bash

linux修改方法：
第一步 删除原有配置
sudo service docker stop
sudo ip link set dev docker0 down
sudo brctl delbr docker0
sudo iptables -t nat -F POSTROUTING
第二步 创建新的网桥
sudo brctl addbr docker0
sudo ip addr add 172.16.10.1/24 dev docker0
sudo ip link set dev docker0 up
第三步 配置Docker的文件
注意： 这里是 增加下面的配置
vi /etc/docker/daemon.json##追加下面的配置即可
{
    "registry-mirrors": ["https://registry.docker-cn.com"],
    "bip": "172.16.10.1/24"
}

systemctl  restart  docker
```
二、改变网络模式，与宿主机共享一个网卡
启动时添加 --net=host

## 实际原因
docker容器开发web程序外部不能访问

最近开发中遇到了一个问题，我使用Dockerfile生成web应用的镜像，在docker容器中运行，测试时发现在外部客户端发起http请求后，cURL返回了错误，error buffer是：Empty reply from server。（本来在本地测一直都是正常的。）说明是外部无法访问这个url。

我排查了很多原因，终于找到是，程序运行的ip写成了app.run(host='127.0.0.1', port=13031)。
改成app.run(host='0.0.0.0', port=13031)就可以正常访问了。

0.0.0.0，localhost和127.0.0.1的区别
在服务器中，0.0.0.0指的是本机上的所有IPV4地址，是真正表示“本网络中的本机”。 一般我们在服务端绑定端口的时候可以选择绑定到0.0.0.0，这样我的服务访问方就可以通过我的多个ip地址访问我的服务。
在路由中，0.0.0.0表示的是默认路由，即当路由表中没有找到完全匹配的路由的时候所对应的路由。
而127.0.0.1是本地回环地址中的一个，大多数windows和Linux电脑上都将localhost指向了127.0.0.1这个地址，相当于是本机地址。
localhost是一个域名，可以用它来获取运行在本机上的网络服务。
在大多数系统中，localhost被指向了IPV4的127.0.0.1和IPV6的::1。

# 关于Ubuntu拒绝root用户ssh远程登录
sudo vim /etc/ssh/sshd_config

找到并用#注释掉这行：PermitRootLogin prohibit-password

新建一行 添加：PermitRootLogin yes

重启服务

sudo service ssh restart

# Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?
wsl -> wsl2
wsl --list --verbose
wsl --set-version Ubuntu-20.04 2
wsl --set-default-version 2

# k8s service 本地能访问，公网不能访问
externalTrafficPolicy: Local -> externalTrafficPolicy: Cluster

# k8s pod 内部ping不通 servicename clusterIP
kubectl edit cm kube-proxy -n kube-system
mode 改为 ipvs
kubectl get pod -n kube-system | grep kube-proxy |awk '{system("kubectl delete pod "$1" -n kube-system")}'

# k8s pod 内部自身服务名访问不通
curl apisix-prometheus.ingress-apisix:9091/apisix/prometheus/metrics
https://github.com/kubernetes/minikube/issues/1568

# listen tcp :8080: bind: An attempt was made to access a socket in a way forbidden by its access permissions.
netstat -aon|findstr "8080"
taskkill /f /pid 12732

# 监听端口外部访问不了,跟docker内监听外部访问不了一样
127.0.0.1 -> 0.0.0.0

# ERROR:  database "lyy" is being accessed by other users DETAIL:  There is 1 other session using the database.
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'timepill';

# es ocker ElasticsearchException[failed to bind service]; nested: AccessDeniedException[/usr/share/elasticsearch/data/nodes];
Dockerfile
FROM docker.elastic.co/elasticsearch/elasticsearch:5.0.0


USER root
RUN chown elasticsearch:elasticsearch -R /usr/share/elasticsearch/data

USER elasticsearch
EXPOSE 9200 9300

暴力
chmod 777 -R /data/es

# Docker挂载的文件(docker run-v)在宿主机修改了后，在容器中没有生效的解决办法
docker run -v 挂载到容器中的文件（注意不是目录）一般是配置文件，在宿主机vi wq之后，进容器里面看发现改动没有生效，后来找了很久没有发现解决办法，直到看到这篇里面提到了需要修改那个文件的权限为666（chmod 666 xxx.conf），但是值得注意的是：中途修改的无效，需要run之前就修改了。

# Drone-Server level=fatal msg="main: invalid configuration" error="Invalid port configuration. See https://discourse.drone.io/t/drone-server-changing-ports-protocol/4144"
这个问题是 k8s 与 drone 之间的命名问题, 官方竟然一直不解决或者明确说明。

解决方式一: 在创建 deployment、StatefulSet、service 不能创建名字为 drone-server 的服务。

解决方式二: 配置 DRONE_SERVER_PORT=:80 变量

# Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running
怎么重启docker都没用
reboot

# k8s etcd集群无法重启，
原因有绑定的 PersistentVolumeClaim

kubectl scale --replicas=0 StatefulSet/apisix-etcd -n ingress-apisix
kubectl delete PersistentVolumeClaim $(kubectl get PersistentVolumeClaim -n ingress-apisix | awk '{print $1}') -n ingress-apisix
kubectl scale --replicas=1 StatefulSet/apisix-etcd -n ingress-apisix
## minikube 有秘钥 /var/lib/minikube/certs/etcd

cp -r /var/lib/minikube/certs/etcd /root/certs &&  chmod 666 /root/certs/etcd/server.key || k8s.runAsUser=0 || initContainer.command - chown -R nobody:nobody /certs/etcd

# Docker 启动alpine镜像中可执行程序文件遇到 not found
问题： docker alpine镜像中遇到 sh: xxx: not found
例如：
在容器内/app/目录下放置了可执行文件abc，启动时提示not found
/app/startup.sh: line 5: ./abc : not found
原因
由于alpine镜像使用的是musl libc而不是gnu libc，/lib64/ 是不存在的。但他们是兼容的，可以创建个软连接过去试试!
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# k8s找不到执行文件
command:
- autossh -M 0 -o StrictHostKeyChecking=no -o ServerAliveInterval=120 -o ServerAliveCountMax=3 -o ConnectTimeout=60 -o ExitOnForwardFailure=yes -CTN -D 0.0.0.0:1080 root@host
应该用
---
command:
- autossh
args:
- "-M 0"
- "-o StrictHostKeyChecking=no"
- "-o ServerAliveInterval=120"
- "-o ServerAliveCountMax=3"
- "-o ConnectTimeout=60"
- "-o ExitOnForwardFailure=yes"
- -CTN
- "-D 0.0.0.0:1080"
- root@$SSH_HOST

# job=\"kubelet\", metrics_path=\"/metrics\", namespace=\"kube-system\", node=\"vm-20-12-ubuntu\", service=\"kube-prometheus-kube-prome-kubelet\"}];many-to-many matching not allowed: matching labels must be unique on one side

# drone \ 反斜杠，sed用的反斜杠，yaml found unknown escape character
解决不了，# 分隔符

# docker遇到的错误:Get https://registry-1.docker.io/v2/: dial tcp: lookup registry-1.docker.io on [::1]:53
vi /etc/resolv.conf
加入：

nameserver y.y.y.y
nameserver x.x.x.x
## wsl没法上网 ping不通 Temporary failure in name resolution
wsl.exe --shutdown
重启

ipconfig
Ethernet adapter vEthernet (WSL):
IPv4 Address. . . . . . . . . . . : 172.17.0.1
猜测ip跟docker0 冲突
重装解决

# docker protoc 
Could not make proto path relative: /work/proto/utils/*.proto: No such file or directory

# docker build ADD failed: file not found in build context or excluded by .dockerignore: file does not exist
所以我在。目录并执行sudo docker构建-t测试7/3/web，它构建容器，但在ADD etc/php /usr/local/etc/php...

您提供给docker build(7/3/web)定义了“构建上下文”。Docker只能访问此目录下的文件。

最简单的解决方案可能是将您的Dockerfile进入toplevel目录，并将您的路径调整为从那里开始的相对路径。

你不需要有要移动Dockerfile (正如Jason所指出的，您可以通过以下方式在任何地方引用它-f)，但这使得目录的组织更加明显。
## cd topdir && docker build -t tag -f
## docker build -t tag -f Dockerfile topdir

# k8s dial tcp: lookup postgres.tools on x.x.x.x: no such host
dnsPolicy: Default #无法解析svc

# k8s docker postgres psql Peer authentication failed for user "postgres"
-h 0.0.0.0

# fcntl64: symbol not found
c - 如何强制链接到较旧的libc`fcntl`而不是`fcntl64`？


似乎GLIBC 2.28 (released August 2018)对fcntl进行了相当激进的更改。在<fcntl.h>中将定义更改为不再是外部函数but a #define to fcntl64。

结果是，如果您使用此glibc在系统上编译代码（如果它完全使用fcntl（）），那么从2018年8月开始，生成的二进制文件将不会在系统上执行。这会影响很多应用程序。 .fcntl（）的手册页显示，这是一小部分子功能的入口：

https://linux.die.net/man/2/fcntl

如果您可以告诉链接器所需的GLIBC函数的特定版本，那就太好了。但是我发现最接近的是在另一篇文章的答案中描述的这个技巧：

Answer to "Linking against older symbol version in a .so file"

这有点复杂。 fcntl是可变参数，而没有接受va_list。在这种情况下，you cannot forward an invocation of a variadic function。 :-(

当一个程序具有稳定的代码且具有较低的依赖关系时，就很难在当前的Ubuntu上构建它了……然后让可执行文件拒绝在仅一年前（即一天）发布的另一个Ubuntu上运行。一个人有什么追索权？

GLIBC没有办法#define USE_FCNTL_NOT_FCNTL64的事实说明了很多。不管是对是错，大多数OS +工具链制造商似乎已经决定，从较新版本中针对较旧版本系统的二进制文件作为目标并非是当务之急。

阻力最小的途径是使虚拟机远离用于构建项目的最旧的OS +工具链。只要您认为二进制文件将在旧系统上运行，就可以使用该文件来生成二进制文件。

但...


如果您认为用法位于fcntl（）调用的子集中，而该子集不受偏移量大小更改的影响（也就是说，您不使用字节范围锁）
或愿意审查偏移量情况的代码，以使用向后兼容的结构定义
并且不害怕伏都教

# frolvlad/alpine-glibc镜像 无法运行go程序 fcntl64: symbol not found
glic升级问题，用alpine镜像 静态编译

# kubernetes svc NodePort设置externalTrafficPolicy:Local无法访问题
kube-proxy proxy-mode ipvs

# minikube Failed to save config: failed to acquire lock for /root/.minikube/profiles/minikube/config.json: unable to open /tmp/juju-mk270d1b5db5965f2dc9e9e25770a63417031943: permission denied
如果因为种种原因，上一步运行的时候报错，那再次执行上一步 操作之前，需要先进行 sudo rm -rf /tmp/juju-mk* sudo rm -rf /tmp/minikube.* 删除操作，否则会报出如下错误：Failed to save config: failed to acquire lock for /root/.minikube/profiles/minikube/config.json: unable to open /tmp/juju-mk270d1b5db5965f2dc9e9e25770a63417031943: permission denied

# 卸载postgresql
apt remove postgresql*

# minikube start 启动不成功
--extra-config=apiserver.service-node-port-range=1-63353 端口号改小点 39999

# nginx openresty failed to connect: no resolver defined to resolve
最近一直在研究 openResty, 使用过程中在用 lua 脚本连接 redis 的时候，使用了阿里云的云 redis，大家都知道的阿里云的云 redis，连接地址是一个域名，这个时候报错 failed to connect: no resolver defined to resolve，先去检查了一下 redis 的白名单，发现内网的 ecsIP 是在白名单的，然后使用 php 测试连接都是正常的，后面去网上查找资料，终于在墙外找到了答案：
nginx 自己的 resolver 目前尚不支持本地的 /etc/hosts文件（注意，这与 DNS 服务本身无关），而 ngx_lua 的 cosocket 也使用的是 nginx 自己的非阻塞的 DNS resolver 组件。所以我们 只需要在 nginx.conf 中加一行:

http 块加 resolver 8.8.8.8; 
resolver 10.96.0.10; 

# k8s dns 
nameserver 10.96.0.10
search default.svc.cluster.local svc.cluster.local cluster.local
options ndots:5
默认域名 svcname.namespace.svc.cluster.local

# Error from server (InternalError): Internal error occurred: Authorization error (user=kube-apiserver-kubelet-client, verb=get, resource=nodes, subresource=proxy)
kubelet bootstrap 引导出错导致kube-apiserve 和 kubelet 之前自动证书审批未完成，导致两者之间未建立连接
删除kubelet证书并重启kubelet 让 kubelet bootstrap 重新引导完成自动证书审批工作 ，问题解决

# --net=host无效，只能用端口映射
经测试是那台服务器问题，具体原因就不知道了，无法监听ipv4
在Windows或MacOS下运行Docker时，实际上是在Linux虚拟机中运行Docker。设置network_mode: host时，您会将容器连接到虚拟机的网络环境。这很少有用。
实际上，只有在Linux上本地运行Docker时，network_mode: host才有意义。
您最好坚持使用端口发布。
同理是这样？也是服务器跑在vm里的
wls2完全没问题啊？？？
curl localhost 127.0.0.1 0.0.0.0 没问题，解决不掉
反而只有docker的端口映射能成功! docker-proxy -proto tcp -host-ip 0.0.0.0 -host-port 8389 -container-ip 172.17.0.2 -container-port 8388
不知原因，无法解决 

# xshell vim 数字键变字母，复制进去乱码
vi是正常的

# apisix路由设置apisix.dev ，配了hosts 总是自动跳转https 307 HSTS 还以为是apisix的问题换了apisix.d可以了
postman是可以的，感觉跟浏览器有关，请求根本没发出去

# frp直接代理https失败
配置改成正确ip就可以了
## kube nodePort 显示监听0.0.0.0 只能通过minikube ip即本机ip访问，不能通过localhost 和 127.0.0.1和0.0.0.0访问
好像只能通过node ip访问
端口可以监听到不同ip

# minikube 对外开放
--apiserver-ips=0.0.0.0（无效）
--extra-config=apiserver.advertise-address=0.0.0.0（无效）

最后发现配错规则了，不是这个机子...