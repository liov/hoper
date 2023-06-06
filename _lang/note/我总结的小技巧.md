# pogres表移动到另一个库
pg_dump -t table_to_copy source_db | psql target_db

pg_dump -U postgres -d test | psql -d hoper -U postgres

# win10教育版
slmgr /ipk NW6C2-QMPVW-D7KKK-3GKT6-VCFB2

slmgr /skms kms.03k.org

slmgr /ato


# CSDN 复制
javascript:document.body.contentEditable='true';document.designMode='on'; void 0

javascript:document.querySelectorAll(".prism").forEach((b)=>{b.onclick = function(e){mdcp.copyCode(e)}});
document.querySelectorAll("style").forEach((s)=>{if((s.innerText||"").indexOf('#content_views pre')>-1){s.parentElement.removeChild(s)}});

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


# cp复制带.git 的目录
git clone /src /dst

# runAsUser
spec:
securityContext:
runAsUser: 0
containers:

# elastic改密码
POST /_security/user/<user>/_password
{
"password" : "new-password"
}

# python

https://mirrors.aliyun.com/pypi/simple/     # 阿里云
https://pypi.douban.com/simple/             # 豆瓣
https://pypi.tuna.tsinghua.edu.cn/simple    # 清华大学

pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
pip install -i  https://pypi.tuna.tsinghua.edu.cn/simple face_recognition
# ubuntu
deb https://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse

# wsl
-- pre
```bash
[boot]
command="/usr/libexec/wsl-systemd"
```
-- now
wsl.conf 的配置设置
wsl.conf 文件基于每个分发配置设置。 (有关 WSL 2 分发版的全局配置，请参阅 .wslconfig) 。

wsl.conf 文件支持四个部分：automount、network和interopuser。 (在.ini文件约定之后建模，密钥将在节下声明，如 .gitconfig files.) 有关存储 wsl.conf 文件的位置的信息，请参阅 wsl.conf 。

systemd 支持
默认情况下，许多 Linux 分发版运行“systemd” (，包括 Ubuntu) 和 WSL 最近添加了对此系统/服务管理器的支持，以便 WSL 更类似于在裸机计算机上使用你喜欢的 Linux 分发版。 需要版本 0.67.6+ 的 WSL 才能启用系统化。 使用命令 wsl --version检查 WSL 版本。 如果需要更新，可以在 Microsoft Store 中获取最新版本的 WSL。 在 博客公告中了解详细信息。

若要启用 systemd，请使用sudo管理员权限在文本编辑器中打开文件wsl.conf，并将以下行添加到/etc/wsl.conf：

```bash
[boot]
systemd=true
```
然后，需要使用 PowerShell 关闭 WSL 分发 wsl.exe --shutdown 版来重启 WSL 实例。 分发重启后，系统应运行。 可以使用以下命令进行确认： systemctl list-unit-files --type=service这将显示服务的状态。

# 数据库尽量不要用bool表示状态
否则你将只能表示两种状态，用int2利于扩展

# Android Emulator Networking
10.0.2.1	Router/gateway address
10.0.2.2	Special alias to your host loopback interface (i.e., 127.0.0.1 on your development machine)
10.0.2.3	First DNS server
10.0.2.4 / 10.0.2.5 / 10.0.2.6	Optional second, third and fourth DNS server (if any)
10.0.2.15	The emulated device's own network/ethernet interface
127.0.0.1	The emulated device's own loopback interface

# go用channel控制批量任务的结束，要用close(ch) 不要用ch<-,只会停掉一个任务，其他会卡住

# linux硬盘清理
/var/log/journal 占用太大
保留一周
journalctl --vacuum-time=1w

保留一月（推荐）
journalctl --vacuum-time=1month

保留一年
journalctl --vacuum-time=1years

保留500M
journalctl --vacuum-size=500M

保留1G
journalctl --vacuum-size=1G

# docker 清理
docker system prune
## 日志
{
"log-driver": "json-file",
"log-opts": {"max-size": "10m", "max-file": "3"}
}

ls -lh $(find /var/lib/docker/containers/ -name *-json.log)

docker run --rm -v /var/lib/docker:/var/lib/docker alpine sh -c "echo '' > $(docker inspect --format='{{.LogPath}}' CONTAINER_NAME)"


truncate -s 0 /var/lib/docker/containers/*/*-json.log
sudo truncate -s 0 `docker inspect --format='{{.LogPath}}' <container>`

# k8s hosts
hostAliases
```yaml
apiVersion: v1
kind: Pod
spec:
  restartPolicy: Never
  hostAliases:
    - ip: "4.1.2.3"
      hostnames:
      - "a.com"
      - "b.com"
```