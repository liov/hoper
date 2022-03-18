
# zsh
```bash
sudo apt install zsh
sh -c "$(curl -fsSL https://raw.github.com/robbyrussell/oh-my-zsh/master/tools/install.sh)"
vi ~/.zshrc
mh
source
erport
```

# [protoc](https://github.com/protocolbuffers/protobuf/releases)


# rust
curl https://sh.rustup.rs -sSf | sh

# [chocolatey](https://chocolatey.org)
```bash
ChocolateyInstall = xxxx/Chocolatey
cmd:
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
ps:
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco upgrade chocolatey
```


# [minikube](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
minikube的ssh有问题，不能发送esc
- https://github.com/kubernetes/minikube/releases/latest/download/minikube-installer.exe
```bash
win:
choco install minikube
unix:
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube
sudo install minikube /usr/local/bin

minikube start --registry-mirror=https://registry.docker-cn.com --vm-driver=hyperv --hyperv-virtual-switch "Default Switch"
```


# bazel[bazel](https://bazel.build,https://github.com/bazelbuild/bazel/releases)


# yarn
wget https://yarnpkg.com/latest.tar.gz

tar -zvxf latest.tar.gz


# git
git config --global user.name ${username}
git config --global user.email ${email}
--global credential.helper store

# ssh

powershell ssh
```
Add-WindowsCapability -Online -Name OpenSSH-Client
  ssh root@IP  -p PORT -i .\.ssh\id_rsa
ssh -R [local port]:[remote host]:[remote port] [SSH hostname]
ssh  -fNg -L <本地端口>:<服务器数据库地址>  <用户名>@<服务器地址>
想让SSH连接一直连接，可以加上 -NTf 参数。
exit
```

# sublimet text
```json
{
    "default_line_ending": "unix",
    "hot_exit": false,
    "remember_open_files": false,
    "theme": "Adaptive.sublime-theme",
    "enable_tab_scrolling": false
}
```

# [go](https://golang.google.cn/doc/)
wget https://dl.google.com/go/go1.xx.x.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.xx.x.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
export GOPRIVATE=go.hoper.xyz
export GOPROXY=https://goproxy.io,direct

# gradle maven
```groovy
maven { url 'https://maven.aliyun.com/repository/public' }//central和jcenter的聚合
maven { url 'https://maven.aliyun.com/repository/central' }
maven { url 'https://maven.aliyun.com/repository/jcenter' }
maven { url 'https://maven.aliyun.com/repository/google'}
maven { url 'https://maven.aliyun.com/repository/gradle-plugin' }
maven { url "https://jitpack.io" }
```

# etcd
etcdctl --endpoints=https://127.0.0.1:2379 --cacert=/etc/kubernetes/pki/etcd/ca.crt --cert=/etc/kubernetes/pki/etcd/peer.crt --key=/etc/kubernetes/pki/etcd/peer.key member list


# IDEA
plugin 仓库地址 https://repo.idechajian.com https://plugins.zhile.io


# [helm](https://helm.sh/)
wget https://get.helm.sh/helm-v3.6.0-linux-amd64.tar.gz
tar -zxvf helm-v3.0.0-linux-amd64.tar.gz
mv linux-amd64/helm /usr/local/bin/helm
---------------------------------------
curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3
chmod 700 get_helm.sh
./get_helm.sh

helm repo add gitlab https://charts.gitlab.io/
helm repo add aliyun https://kubernetes.oss-cn-hangzhou.aliyuncs.com/charts
helm repo add incubator https://kubernetes-charts-incubator.storage.googleapis.com/


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

## 比较骚
1.操作位置：迁移数据库源（旧数据库主机）

找到PostgreSql 的data目录   关闭数据库进程

打包 tar -zcvf pgdatabak.tar.gz data/

------------------------------------------------------------------

2.通过winScp 或者 CRT 等工具拷贝到    迁移目标源（新主机--需安装postgresql）  同样的data目录 关闭数据库进程

解压  tar -zxvf pgdatabak.tar.gz -C /usr/local/postgres/

重新授权 执行命令  chown -R postgres.postgres data/