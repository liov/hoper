sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 3B4FE6ACC0B21F32 871920D1991BC93C 7EA0A9C3F273FCD8

sudo vi /etc/apt/sources.list

2.1 阿里源：
deb http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-security main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-updates main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src http://mirrors.aliyun.com/ubuntu/ bionic-backports main restricted universe multiverse
2.2 其它源


#清华源
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb-src https://mirrors.tuna.tsinghua.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
 
#中科大源
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb https://mirrors.ustc.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-updates main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-backports main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-security main restricted universe multiverse
deb-src https://mirrors.ustc.edu.cn/ubuntu/ bionic-proposed main restricted universe multiverse
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
- WSL2使用的是虚拟机技术和WSL第一版本不一样，和宿主windows不在同一个网络内
- 获取宿主windows的ip
export windows_host=`ipconfig.exe | grep -n4 WSL  | tail -n 1 | awk -F":" '{ print $2 }'  | sed 's/^[ \r\n\t]*//;s/[ \r\n\t]*$//'`

- 假设你的宿主windows代理端口是1080, 全面设置WSL内的代理
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
                         
# 设置docker内的代理
在/etc/default/docker有
 export http_proxy="socks5://[windows_host]:1080"
 export https_proxy="socks5://[windows_host]:1080"
sudo sed -i -E "s#socks5.*?1080#socks5://$windows_host:1080#" /etc/default/docker

# 设置WSL内的DNS, 默认是系统自己创建的
sudo bash -c 'echo -e "\nnameserver 114.114.114.114\nnameserver 8.8.8.8\nnameserver 8.8.4.4" > /etc/resolv.conf'
```
## Note, this file is written by cloud-init on first boot of an instance
## modifications made here will not survive a re-bundle.
## if you wish to make changes you can:
## a.) add 'apt_preserve_sources_list: true' to /etc/cloud/cloud.cfg
##     or do the same in user-data
## b.) add sources in /etc/apt/sources.list.d
## c.) make changes to template file /etc/cloud/templates/sources.list.tmpl

# See http://help.ubuntu.com/community/UpgradeNotes for how to upgrade to
# newer versions of the distribution.
deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal main
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal main

## Major bug fix updates produced after the final release of the
## distribution.
deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates main
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates main

## N.B. software from this repository is ENTIRELY UNSUPPORTED by the Ubuntu
## team. Also, please note that software in universe WILL NOT receive any
## review or updates from the Ubuntu security team.
deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal universe
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal universe
deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates universe
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates universe

## N.B. software from this repository is ENTIRELY UNSUPPORTED by the Ubuntu
## team, and may not be under a free licence. Please satisfy yourself as to
## your rights to use the software. Also, please note that software in
## multiverse WILL NOT receive any review or updates from the Ubuntu
## security team.
# deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal multiverse
# deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal multiverse
# deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates multiverse
# deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-updates multiverse

## Uncomment the following two lines to add software from the 'backports'
## repository.
## N.B. software from this repository may not have been tested as
## extensively as that contained in the main release, although it includes
## newer versions of some applications which may provide useful features.
## Also, please note that software in backports WILL NOT receive any review
## or updates from the Ubuntu security team.
# deb http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-backports main restricted universe multiverse
# deb-src http://mirrors.cloud.aliyuncs.com/ubuntu/ focal-backports main restricted universe multiverse

## Uncomment the following two lines to add software from Canonical's
## 'partner' repository.
## This software is not part of Ubuntu, but is offered by Canonical and the
## respective vendors as a service to Ubuntu users.
# deb http://archive.canonical.com/ubuntu focal partner
# deb-src http://archive.canonical.com/ubuntu focal partner

deb http://mirrors.cloud.aliyuncs.com/ubuntu focal-security main
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu focal-security main
deb http://mirrors.cloud.aliyuncs.com/ubuntu focal-security universe
deb-src http://mirrors.cloud.aliyuncs.com/ubuntu focal-security universe
# deb http://mirrors.cloud.aliyuncs.com/ubuntu focal-security multiverse
# deb-src http://mirrors.cloud.aliyuncs.com/ubuntu focal-security multiverse
deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable
# deb-src [arch=amd64] https://download.docker.com/linux/ubuntu focal stable


deb http://mirrors.tencentyun.com/ubuntu/ focal main restricted universe multiverse
deb http://mirrors.tencentyun.com/ubuntu/ focal-security main restricted universe multiverse
deb http://mirrors.tencentyun.com/ubuntu/ focal-updates main restricted universe multiverse
# deb http://mirrors.tencentyun.com/ubuntu/ focal-proposed main restricted universe multiverse
# deb http://mirrors.tencentyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb-src http://mirrors.tencentyun.com/ubuntu/ focal main restricted universe multiverse
deb-src http://mirrors.tencentyun.com/ubuntu/ focal-security main restricted universe multiverse
deb-src http://mirrors.tencentyun.com/ubuntu/ focal-updates main restricted universe multiverse
# deb-src http://mirrors.tencentyun.com/ubuntu/ focal-proposed main restricted universe multiverse
# deb-src http://mirrors.tencentyun.com/ubuntu/ focal-backports main restricted universe multiverse
deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable
# deb-src [arch=amd64] https://download.docker.com/linux/ubuntu focal stable

echo 'deb http://mirrors.cloud.tencent.com/debian/ buster main non-free contrib\ndeb http://mirrors.cloud.tencent.com/debian-security buster/updates main\ndeb http://mirrors.cloud.tencent.com/debian/ buster-updates main non-free contrib\ndeb http://mirrors.cloud.tencent.com/debian/ buster-backports main non-free contrib' > /etc/apt/sources.list
