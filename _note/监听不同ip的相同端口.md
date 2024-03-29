## 监听总结之证明 IP地址各自拥有自己的端口号空间
IP地址各自拥有自己的端口号空间，比如同一个主机上有三个ip地址，其中127.0.0.1、192.168.0.10和192.168.0.100，而192.168.0.100是附加绑定在192.168.0.10所在的网卡上，也即192.168.0.100和192.168.0.10在同一个网卡上。“IP地址各自拥有自己的端口号空间”。

127.0.0.1有域名localhost
## TCP和UDP可以同时监听相同的端口号吗？
同时监听同一端口号的理解：同一时刻,TCP的某一端口和UDP的相同数字端口都是listening状态。
原因：1.端口并不是物理概念，而仅仅只是两个字节。
2.UDP的端口和TCP的端口是两个东西，互不影响。
3.TCP和UDP协议监听同一个端口后，接收到的数据并不会产生影响。因为：接收数据是按照 五元组{传输协议，源IP，目的IP，源端口号，目的端口号}来判断接收者。

## 监听不同ip地址的差异
后端监听ip地址的三种主要的方式
- 监听到127.0.0.1
- 监听到0.0.0.0
- 监听到主机内网ip

### 监听127.0.0.1
本机通过127.0.0.1访问成功，网络接口为loopback
本机通过局域网IP 192.168.0.113访问失败，网络接口为loopback
同一局域网下的外部主句通过局域网IP 192.168.0.112访问失败，网络接口-et1
因此，在实际应用中，我们在服务端监听ip地址的时候不要绑定到127.0.0.1，如果绑定到了127.0.0.1，会导致我们的应用只能在本地127.0.0.1访问，其他人无法通过其他任何方式进行访问

### 监听0.0.0.0
本机通过127.0.0.1访问成功，网络接口为loopback
本机通过局域网IP 192.168.0.113 访问成功，网络接口为loopback
同一局域网下的外部主句通过局域网IP 192.168.0.112访问成功，网络接口-et1
比如我有一台服务器，一个外网A,一个内网B，如果我绑定的端口指定了0.0.0.0，那么通过内网地址或外网地址都可以访问我的应用。

### 监听主机内网ip 192.168.0.113
本机通过127.0.0.1访问失败，网络接口为loopback
本机通过局域网IP 192.168.0.113 访问成功，网络接口为loopback
同一局域网下的外部主句通过局域网IP 192.168.0.112访问成功，网络接口-et1

# 总结
在实际应用中，最好的监听ip地址方式为：监听到0.0.0.0