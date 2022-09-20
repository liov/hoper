# step 0: clear# sudo apt-get remove docker docker-engine docker-ce docker.io# step 1: 安装必要的一些系统工具sudo apt-get update#sudo apt-get -y install apt-transport-https ca-certificates curl software-properties-common
# step 2: 安装GPG证书# curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -    //官方
# curl -fsSL https://mirrors.aliyun.com/docker-ce/linux/ubuntu/gpg | sudo apt-key add -   //阿里云
# Step 3: 写入软件源信息# sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"  //官方
# sudo add-apt-repository "deb [arch=amd64] https://mirrors.aliyun.com/docker-ce/linux/ubuntu $(lsb_release -cs) stable"  //阿里云
# Step 4: 更新并安装Docker-CE
# sudo apt-get -y update
# sudo apt-get -y install docker-ce# 此处默认安装的是docker最新版本，但由于WSL的问题，可能会存在问题，会在拉取镜像的时候出现错误：Error response from daemon: OCI runtime create failed context canceled# 所以我们安装指定的docker版本，见4.1，我安装的是18.03.1-ce

# Step 4.1 安装指定版本的Docker-CE:
# Step 4.11: 查找Docker-CE的版本:
# apt-cache madison docker-ce
#   docker-ce | 17.03.1~ce-0~ubuntu-xenial | https://mirrors.aliyun.com/docker-ce/linux/ubuntu xenial/stable amd64 Packages
#   docker-ce | 17.03.0~ce-0~ubuntu-xenial | https://mirrors.aliyun.com/docker-ce/linux/ubuntu xenial/stable amd64 Packages
# Step 4.1.2: 安装指定版本的Docker-CE: (VERSION例如上面的17.03.1~ce-0~ubuntu-xenial)
# sudo apt-get -y install docker-ce=[VERSION]# Step 5: 查看docker的版本# docker -v# Step 6: 查看docker的启动状态# sudo service docker status# Step 7：启动docker# sudo service docker start  //备注：WSL目前是WSL1，所以linux的发行版本部分命令还不能使用，比如systemctl命令。此处启动docker可以使用：sudo systemctl start docker# Step 8：拉取hello-world镜像，验证docker是否正常# docker run hello-world# 输出：# Hello from Docker!# This message shows that your installation appears to be working correctly.# ......

sudo gpasswd -a ${USER} docker
newgrp docker

echo "export DOCKER_HOST=tcp://192.168.1.212:2375" >> ~/.bashrc && source ~/.bashrc


1.打开编辑：
vi /lib/systemd/system/docker.service

2.注释原有的：
#ExecStart=/usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock

3.添加新的：
ExecStart=/usr/bin/dockerd -H unix:///var/run/docker.sock -H tcp://0.0.0.0:2375

-H代表指定docker的监听方式，这里是socket文件文件位置，也就是socket方式，2375就是tcp端口
1
4.保存并退出

5.重新加载系统服务配置文件（包含刚刚修改的文件）
systemctl daemon-reload

6.重启docker服务
systemctl restart docker
