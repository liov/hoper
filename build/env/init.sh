zjokwzwojyfvbgef

ps -aux | grep main

#使用命令：

#netstat -apn

#kill -9 [PID]
git checkout init.sh
git checkout ../../micro/client/hoper
git pull

#git status
#git checkout


sudo vi /etc/apt/sources.list

# 默认注释了源码镜像以提高 apt update 速度，如有需要可自行取消注释
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse
#deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse
# 预发布软件源，不建议启用
# deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-proposed main restricted universe multiverse
# deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-proposed main restricted universe multiverse

sudo apt-get update
sudo apt-get upgrade



、添加新的用户账号使用 useradd命令，其语法如下：
代码:
useradd 选项 用户名
其中各选项含义如下：
代码:
-c comment 指定一段注释性描述。
-d 目录 指定用户主目录，如果此目录不存在，则同时使用-m选项，可以创建主目录。
-g 用户组 指定用户所属的用户组。
-G 用户组，用户组 指定用户所属的附加组。
-s Shell文件 指定用户的登录Shell。
-u 用户号 指定用户的用户号，如果同时有-o选项，则可以重复使用其他用户的标识号。
用户名 指定新账号的登录名。
例1：
代码:
# useradd –d /usr/sam -m sam
此命令创建了一个用户sam，
其中-d和-m选项用来为登录名sam产生一个主目录/usr/sam（/usr为默认的用户主目录所在的父目录）。
例2：
代码:
# useradd -s /bin/sh -g group –G adm,root gem
此命令新建了一个用户gem，该用户的登录Shell是/bin/sh，它属于group用户组，同时又属于adm和root用户组，其中group用户组是其主组。
这里可能新建组：#groupadd group及groupadd adm　
增加用户账号就是在/etc/passwd文件中为新用户增加一条记录，同时更新其他系统文件如/etc/shadow, /etc/group等。
Linux提供了集成的系统管理工具userconf，它可以用来对用户账号进行统一管理。

useradd -d /home/jyb -m -g root jyb

sudo apt update

sudo apt intall git

wget https://ftp.gnu.org/gnu/gcc/gcc-8.2.0/gcc-8.2.0.tar.gz
wget https://ftp.gnu.org/gnu/gmp/gmp-4.3.2.tar.gz
wget https://ftp.gnu.org/gnu/mpfr/mpfr-2.4.2.tar.gz
wget https://ftp.gnu.org/gnu/mpc/mpc-1.0.1.tar.gz

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

sudo apt intall screen

wget https://yarnpkg.com/latest.tar.gz

tar zvxf latest.tar.gz
