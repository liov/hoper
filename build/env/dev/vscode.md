vscode的启动速度已经快到可以接受的程度了
似乎要与sublime说再见了

remote-ssh 设置配置文件

ssh $user@$domain -p $port -A

Host dev
    HostName 192.168.1.212
    User crm
    Port 10000
    ForwardAgent no #ssh转发
    IdentityFile ~/.ssh/id_rsa-remote-ssh
    DynamicForward 1080
    ForwardX11 yes
    ForwardX11Trusted yes
    
连接输密码（推荐使用IdentityFile ）

vscode remote 爱了啊，打印调用行数能直接跳转,除了不能保存密码
要不是xshell极低的资源占用

```Dockerfile
FROM bitnami/minideb

ENV VSCODE_SERVER_DATA_DIR /data

RUN apt-get update && apt-get install -y wget gnome-keyring

RUN wget -O- https://aka.ms/install-vscode-server/setup.sh | sh

CMD ["code-server", "serve-local" ,"--host","0.0.0.0" , "--accept-server-license-terms"]
```