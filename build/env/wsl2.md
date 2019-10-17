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
