
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

minikube start --registry-mirror= --vm-driver=hyperv --hyperv-virtual-switch "Default Switch"
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


# snap

apt insatall snapd

# linux
vim /etc/profile
export PATH=$PATH:/usr/local/go/bin:/usr/local/node/bin
export HISTCONTROL=ignoredups
source /etc/profile

# netcat
https://eternallybored.org/misc/netcat/

# gcc
wget https://ftp.gnu.org/gnu/gcc/gcc-8.2.0/gcc-8.2.0.tar.gz
wget https://ftp.gnu.org/gnu/gmp/gmp-4.3.2.tar.gz
wget https://ftp.gnu.org/gnu/mpfr/mpfr-2.4.2.tar.gz
wget https://ftp.gnu.org/gnu/mpc/mpc-1.0.1.tar.gz
```bash
mkdir ~/local/gcc

tar xf gmp-4.3.2.tar.gz
cd gmp-4.3.2
sudo yum install -y m4
./configure --prefix=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf mpfr-2.4.2.tar.gz
cd mpfr-2.4.2
./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf mpc-1.0.1.tar.gz
cd mpc-1.0.1
./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc --with-mpfr=$HOME/local/gcc
make && make install
ls ~/local/gcc/lib/

cd ..
tar xf gcc-8.2.0.tar.gz
cd gcc-8.2.0

./configure --prefix=$HOME/local/gcc --with-gmp=$HOME/local/gcc --with-mpfr=$HOME/local/gcc --with-mpc=$HOME/local/gcc --disable-multilib
export LD_LIBRARY_PATH=$HOME/local/gcc/lib:$LD_LIBRARY_PATH
make && make install

export PATH=$HOME/local/gcc/bin:$PATH
```

# vsftp 

```bash
sudo apt install vsftpd

sudo passwd jyb

mkdir /home/jyb/ftp

chmod 777 -R /home/jyb/ftp

sudo vim /etc/vsftpd.conf

connect_from_port_21=YES

local_root=/home/jyb/ftp

allow_writeable_chroot=YES

将#chroot_local_user=YES前的注释去掉

pam_service_name=ftp原配置中为vsftpd，ubuntu用户需要更改成ftp

sudo service vsftpd start

sudo service vsftpd restart

```