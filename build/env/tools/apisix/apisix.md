Configure and Installation
APISIX Installed and tested in the following systems:

CentOS 7, Ubuntu 16.04, Ubuntu 18.04, Debian 9, Debian 10, macOS, ARM64 Ubuntu 18.04

There are several ways to install the Apache Release version of APISIX:

Source code compilation (applicable to all systems)

Installation runtime dependencies: OpenResty and etcd, and compilation dependencies: luarocks. Refer to install dependencies documentation
Download the latest source code release package:
$ mkdir apisix-2.0
$ cd apisix-2.0
$ wget https://downloads.apache.org/apisix/2.0/apache-apisix-2.0-src.tgz
$ tar zxvf apache-apisix-2.0-src.tgz
Install the dependencies：
$ make deps
check version of APISIX:
$ ./bin/apisix version
start APISIX:
$ ./bin/apisix start
Docker image （applicable to all systems）

By default, the latest Apache release package will be pulled:

$ docker pull apache/apisix
The Docker image does not include etcd, you can refer to docker compose example to start a test cluster.

RPM package（only for CentOS 7）

Installation runtime dependencies: OpenResty and etcd, refer to install dependencies documentation
install APISIX：
$ sudo yum install -y https://github.com/apache/apisix/releases/download/2.0/apisix-2.0-0.el7.noarch.rpm
check version of APISIX:
$ apisix version
start APISIX:
$ apisix start
Note: Apache APISIX would not support the v2 protocol of etcd anymore since APISIX v2.0, and the minimum etcd version supported is v3.4.0. Please update etcd when needed. If you need to migrate your data from etcd v2 to v3, please follow etcd migration guide.
