CMake >= 3.8.0
Rust >= 1.19.0
binutils >= 2.22[https://mirrors.tuna.tsinghua.edu.cn/gnu/binutils/]
LLVM and Clang >= 3.9 if you need to generate bindings at compile time.
By default, the secure feature is provided by boringssl, which requires Go (>=1.7) to build. You can also use openssl instead by enabling openssl feature.
For Windows, you also need to install following software:

Active State Perl[http://www.perl.org]
yasm[https://github.com/yasm/yasm/releases]
Visual Studio 2015+

cargo install protobuf-codegen
cargo install grpcio-compiler
protoc --rust_out=./protobuf --grpc_out=./protobuf --plugin=protoc-gen-grpc=`which grpc_rust_plugin` ../../../../proto/helloworld.proto


[wsl2]
sudo nc -l 80
telnet 172.27.175.35 80
netstat -a

[keng]
wsl2和windows的网络问题
监听127.0.0.1和0.0.0.0时，natstat显示的Local Address是不同的，
    localhost:50051 [::]:50051
0.0.0.0 不能ping通，代表本机所有的IP地址；
    监听127.0.0.1，创建Socket，那么用本机地址建立tcp连接不成功，反过来也是如此；也就是，监听时采用的地址为192.168.0.1，就只能用192.168.0.1进行连接。
    而监听0.0.0.0创建Socket，那么无论使用127.0.0.1或本机ip都可以建立tcp连接,也就是不论通过127.0.0.1或192.168.0.1、192.168.1.1都能连接成功。
    0.0.0.0建立tcp连接的时候也可以通过绑定IP_ADDR_ANY来实现。
IPv4 的环回地址是保留地址之一 127.0.0.1。尽管只使用 127.0.0.1 这一个地址，但地址 127.0.0.0 到 127.255.255.255 均予以保留。此地址块中的任何地址都将环回到本地主机中。此地址块中的任何地址都绝不会出现在任何网络中。
首先我们来讲讲127.0.0.1，172.0.0.1是回送地址，localhost是本地DNS解析的127.0.0.1的域名，在hosts文件里可以看到。

一般我们通过ping 127.0.0.1来测试本地网络是否正常。其实从127.0.0.1~127.255.255.255，这整个都是回环地址。这边还要

注意的一点就是localhost在了IPV4的是指127.0.0.1而IPV6是指::1。当我们在服务器搭建了一个web服务器的时候如果我们

监听的端口时127.0.0.1：端口号 的 时候，那么这个web服务器只可以在服务器本地访问了，在别的地方进行访问是不行的。

（127.0.0.1只可以在本地ping自己的，那么你监听这个就只可以在本地访问了）

  然后我们来讲讲0.0.0.0，如果我们直接ping 0.0.0.0是不行的，他在IPV4中表示的是无效的目标地址，但是在服务器端它表示

本机上的所有IPV4地址，如果一个服务有多个IP地址（192.168.1.2和10.1.1.12），那么我们如果设置的监听地址是0.0.0.0那

么我们无论是通过IP192.168.1.2还是10.1.1.12都是可以访问该服务的。在路由中，0.0.0.0表示的是默认路由，即当路由表中

没有找到完全匹配的路由的时候所对应的路由。
