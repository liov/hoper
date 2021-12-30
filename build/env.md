
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


# node
wget https://nodejs.org/dist/v12.3.1/node-v12.3.1.tar.gz
tar -xzvf
./configure
apt install g++
wget https://nodejs.org/dist/v12.3.1/node-v12.3.1-linux-x64.tar.xz

tar -Jxvf


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

# protoc卡住
tags标签写错

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

# [Apache APISIX Helm Chart](https://apisix.apache.org/)
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add apisix https://charts.apiseven.com
helm repo update
helm install apisix apisix/apisix
helm install apisix-dashboard apisix/apisix-dashboard
helm install apisix-ingress-controller apisix/apisix-ingress-controller --namespace ingress-apisix
v1.14-v1.19
helm install apisix-ingress-controller apisix/apisix-ingress-controller --namespace ingress-apisix --set config.kubernetes.ingressVersion=networking/v1beta1
## if etcd export by kubernetes service need spell fully qualified name
$ helm install apisix apisix/apisix \
    --set etcd.enabled=false \
    --set etcd.host={http://etcd_node_1:2379\,http://etcd_node_2:2379} \
    --set admin.allow.ipList="{0.0.0.0/0}" \
    --namespace ingress-apisix

helm install apisix-ingress-controller apisix/apisix-ingress-controller \
  --set image.tag=dev \
  --set config.apisix.baseURL=http://apisix-admin:9180/apisix/admin \
  --set config.apisix.adminKey=edd1c9f034335f136f87ad84b625c8f1 \
  --namespace ingress-apisix

# win10教育版
slmgr /ipk NW6C2-QMPVW-D7KKK-3GKT6-VCFB2

slmgr /skms kms.03k.org

slmgr /ato

# kong gateway
kong.conf
```properties
database = postgres
pg_host = postgre.miz.hk
pg_port = 5432
pg_user = web
pg_password = 123456
pg_database = openmng-gw
admin_listen = 0.0.0.0:8001, 0.0.0.0:8444 ssl
plugins = bundled,session,request-inspector,session-go
lua_package_path = /usr/local/?.lua;/usr/local/?/init.lua;

pluginserver_names = go

pluginserver_go_socket = /usr/local/kong/go_pluginserver.sock
pluginserver_go_start_cmd = /usr/local/bin/go-pluginserver -kong-prefix /usr/local/kong/ -plugins-directory /usr/local/kong/go-plugins
pluginserver_go_query_cmd = /usr/local/bin/go-pluginserver -dump-all-plugins -plugins-directory /usr/local/kong/go-plugins
nginx_user = root
```
```bash
#!/usr/bin/env bash

cd go-plugins/session-go
go build github.com/Kong/go-pluginserver
go build -o session-go.so -buildmode plugin session_validator.go check_path.go config.go exchange_token.go redis.go
cd ../../

cat > Dockerfile <<- EOF
FROM kong:2.5-centos

USER root

RUN rm /etc/localtime
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD ./deploy/${1}/kong/kong.conf /etc/kong/

EOF

for dir in plugins/*
do
if [ -d ${dir} ]
then
    echo "ADD ./${dir} /usr/local/share/lua/5.1/kong/${dir}"  >> Dockerfile
fi
done

echo "ADD ./go-plugins/session-go/go-pluginserver /usr/local/bin/go-pluginserver"  >> Dockerfile
echo "ADD ./go-plugins/session-go/session-go.so /usr/local/kong/go-plugins/"  >> Dockerfile

docker build . -t $2
```

# CSDN 复制
javascript:document.body.contentEditable='true';document.designMode='on'; void 0

# python

pip install paddleocr -i https://mirror.baidu.com/pypi/simple 