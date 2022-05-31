# pogres表移动到另一个库
pg_dump -t table_to_copy source_db | psql target_db

pg_dump -U postgres -d test | psql -d hoper -U postgres

# win10教育版
slmgr /ipk NW6C2-QMPVW-D7KKK-3GKT6-VCFB2

slmgr /skms kms.03k.org

slmgr /ato


# CSDN 复制
javascript:document.body.contentEditable='true';document.designMode='on'; void 0

# python

pip install paddleocr -i https://mirror.baidu.com/pypi/simple

# postgres迁移

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

## 比较骚，只适用同版本
1.操作位置：迁移数据库源（旧数据库主机）

找到PostgreSql 的data目录   关闭数据库进程

打包 tar -zcvf pgdatabak.tar.gz data/

------------------------------------------------------------------

2.通过winScp 或者 CRT 等工具拷贝到    迁移目标源（新主机--需安装postgresql）  同样的data目录 关闭数据库进程

解压  tar -zxvf pgdatabak.tar.gz -C /usr/local/postgres/

重新授权 执行命令  chown -R postgres:postgres data/

# IDEA
plugin 仓库地址 https://repo.idechajian.com https://plugins.zhile.io

# k8s的pod镜像的时区正确设置方法
在k8s中部署pod时，很多时候我们使用的镜像不是我们自己制作的，自己制作的可以把时区设置好，但使用别人的镜像时，这些镜像的时区有可能是UTC，比我们的时间少了8小时，看一些日志时很别扭，比较方便的办法是在部署时设置env环境变量，加上

env:
- name: TZ
  value: Asia/Shanghai
如果是helm的包，特别的bitnami的包，通常都env的扩展参数，但这个参数一般在README中都是[]，那想在helm install --set方式如何设置呢？在网上搜索了几个小时，终于stackoverflow上找到一种可用的方式，这里我以安装postgresql的设置为例：

--set postgresql.extraEnvVars\[0\].name=TZ,postgres
ql.extraEnvVars\[0\].value=Asia\/Shanghai
如果你也遇到这种情况，不防试试，也许可以解决

# 利用docker编译
docker run -v "$GOPATH":/go --rm -v "$PWD":/app -w /app -e GOOS="darwin" -e GOARCH="amd64" golang:1.8 go build -v

# clusterIP: None的即为headless service
type: ClusterIP
clusterIP: None
具体表现service没有自己的虚拟IP,nslookup会出现所有pod的ip.但是ping的时候只会出现第一个pod的ip
service没有负载均衡
检查一下是否用了headless service.headless service是不会自动负载均衡的

# alpine 镜像执行go二进制文件
## 编译静态链接
go build -tags netgo
## alpine-glibc镜像

# cp复制带.git 的目录
git clone /src /dst
