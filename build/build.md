hosts 192.30.253.112	github.com
      199.232.5.194	github.global.ssl.fastly.net
go get -u github.com/gpmgo/gopm
安装protoc[https://github.com/protocolbuffers/protobuf/releases]
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/micro/micro
go get -u github.com/micro/protoc-gen-micro

wls2
```bash
Windows Registry Editor Version 5.00

[HKEY_LOCAL_MACHINE\SOFTWARE\Classes\Directory\background\shell\Bash]
@="Bash here"

[HKEY_LOCAL_MACHINE\SOFTWARE\Classes\Directory\background\shell\Bash\command]
@="C:\\Users\\[your-name]\\AppData\\Local\\Microsoft\\WindowsApps\\wt.exe"

[HKEY_LOCAL_MACHINE\SOFTWARE\Classes\Directory\shell\Bash]
@="Bash here"

[HKEY_LOCAL_MACHINE\SOFTWARE\Classes\Directory\shell\Bash\command]
@="C:\\Users\\[your-name]\\AppData\\Local\\Microsoft\\WindowsApps\\wt.exe"
~/.bashrc
# WSL2使用的是虚拟机技术和WSL第一版本不一样，和宿主windows不在同一个网络内
# 获取宿主windows的ip
export windows_host=`ipconfig.exe | grep -n4 WSL  | tail -n 1 | awk -F":" '{ print $2 }'  | sed 's/^[ \r\n\t]*//;s/[ \r\n\t]*$//'`

# 假设你的宿主windows代理端口是1080, 全面设置WSL内的代理
export ALL_PROXY=socks5://$windows_host:1080
export HTTP_PROXY=$ALL_PROXY
export http_proxy=$ALL_PROXY
export HTTPS_PROXY=$ALL_PROXY
export https_proxy=$ALL_PROXY

# 设置git的代理
if [ "`git config --global --get proxy.https`" != "socks5://$windows_host:1080" ]; then
    git config --global proxy.https socks5://$windows_host:1080
fi

# wcd
# cd C:\\ 自动切换到 /mnt/c
function wcd() {
    command cd `wslpath "$1"`
}
                         
#设置docker内的代理
#在/etc/default/docker有
# export http_proxy="socks5://[windows_host]:1080"
# export https_proxy="socks5://[windows_host]:1080"
sudo sed -i -E "s#socks5.*?1080#socks5://$windows_host:1080#" /etc/default/docker

# 设置WSL内的DNS, 默认是系统自己创建的
sudo bash -c 'echo -e "\nnameserver 114.114.114.114\nnameserver 8.8.8.8\nnameserver 8.8.4.4" > /etc/resolv.conf'
```

安装chocolatey[https://chocolatey.org]
```bash
ChocolateyInstall = xxxx/Chocolatey
cmd:
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command "iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))" && SET "PATH=%PATH%;%ALLUSERSPROFILE%\chocolatey\bin"
ps:
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))
choco upgrade chocolatey
```
安装docker
```bash
sudo apt-get update
sudo apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io
sudo gpasswd -a ${USER} docker
sudo service docker restart
vi /etc/docker/daemon.json
{
    "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn"]
}
```

## 安装minikube[https://kubernetes.io/docs/tasks/tools/install-kubectl/]
minikube的ssh有问题，不能发送esc
- https://github.com/kubernetes/minikube/releases/latest/download/minikube-installer.exe
```shell script
win:
choco install minikube
unix:
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 && chmod +x minikube
sudo install minikube /usr/local/bin
```
## 安装k8s单机集群
```bash
minikube start --registry-mirror=https://registry.docker-cn.com --vm-driver=hyperv --hyperv-virtual-switch "Default Switch"
```
## 安装kubectl[https://kubernetes.io/docs/tasks/tools/install-kubectl/]
```shell script
linux:
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version

ubuntu:
sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

win:
https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/windows/amd64/kubectl.exe
```
## 配置kubectl
```shell script
$HOME/.kube/${cluster}

# 配置一个名为 ${cluster} 的集群，并指定服务地址与根证书
kubectl config set-cluster ${cluster} --server=${cluster-server} --certificate-authority=$HOME/.kube/${cluster}/ca.pem
# 设置一个用户为 ${Username} ，并配置访问的授权文件
kubectl config set-credentials ${Username} --client-certificate=${name}.crt --client-key=${name}.pem --embed-certs=true
# 设置一个名为 ${context} 使用 ${cluster} 集群与 ${Username} 用户的上下文，
kubectl config set-context ${context} --cluster=${cluster} --user=${Username}
# 启用 ${context} 
kubectl config use-context ${context}

${cluster} 集群名称，可自己定义
${cluster-server} 集群地址，取值如下： 
prod: https://hoper.xyz
test: https:/hoper.xyz
dev: https://hoper.xyz
${Username} 用户名称，自行定义,注意不同context的用户名不能相同，否则会覆盖配置
${name} 邮件回复的授权文件名，可自行修改，文件名与配置一致即可
${context} 上下文，自行定义
```
下载istio[https://github.com/istio/istio/releases/tag/1.3.0]
```bash
linux:
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.3.0 sh -
cd istio-1.3.0
export PATH=$PWD/bin:$PATH
```
安装istio[https://istio.io/docs/setup/install/kubernetes/]
```bash
linux:
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done
win ps:
Get-ChildItem e:/istio-1.3.0/install/kubernetes/helm/istio-init/files/ | ForEach-Object -Process{
    if ($_.Name -match "crd*yaml"])
    {
        kubectl apply -f $_;
    }

}

kubectl apply -f install/kubernetes/istio-demo.yaml
kubectl label namespace default istio-injection=enabled
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
```
卸载istio
```bash
kubectl delete -f install/kubernetes/istio-demo.yaml
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl delete -f $i; done
```

部署
```bash
当您使用时部署应用程序时kubectl apply，如果Istio边车注入器 在标有istio-injection=enabled以下标记的名称空间中启动，它们将自动将Envoy容器注入您的应用程序窗格：
kubectl label namespace <namespace> istio-injection=enabled
kubectl create -n <namespace> -f <your-app-spec>.yaml
在没有istio-injection标签的命名空间中，您可以使用 istioctl kube-inject 在部署它们之前在应用程序窗格中手动注入Envoy容器：
istioctl kube-inject -f <your-app-spec>.yaml | kubectl apply -f -
```
安装bazel[https://bazel.build,https://github.com/bazelbuild/bazel/releases]

安装postgresql[https://www.postgresql.org/download/]