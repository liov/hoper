password=password
host=host
mkdir ss
cd ss
cat > ssconfig.json <<- EOF
{
    "server": "0.0.0.0",
    "server_port": 8388,
    "password": "${password}",
    "timeout": 300,
    "plugin": "v2ray-plugin",
    "plugin_opts": "server;tls;host=${host}",
    "method": "aes-256-gcm"
}
EOF

#curl https://get.acme.sh | sh -s email=lby.i@qq.com
#~/.acme.sh/acme.sh --issue --dns dns_cf -d ${host}
#~/.acme.sh/acme.sh --issue --dns -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please
#~/.acme.sh/acme.sh --renew -d ${host} --yes-I-know-dns-manual-mode-enough-go-ahead-please


cd ..
git clone https://github.com/shadowsocks/v2ray-plugin.git
docker pull golang
export GOPATH=~/gopath
docker run -v "$GOPATH":/go --rm -v "$PWD/v2ray-plugin":/app -w /app -e GOOS="linux" -e GOARCH="amd64" golang go build -v
mv v2ray-plugin/v2ray-plugin ss

cd ss
cat > Dockerfile <<- EOF
FROM ghcr.io/shadowsocks/ssserver-rust:latest

USER root

RUN cd /tmp && \
 TAG=$(wget -qO- https://api.github.com/repos/shadowsocks/v2ray-plugin/releases/latest | grep tag_name | cut -d '"' -f4) && \
 wget https://github.com/shadowsocks/v2ray-plugin/releases/download/$TAG/v2ray-plugin-linux-amd64-$TAG.tar.gz && \
 tar -xf *.gz && \
 rm *.gz && \
 mv v2ray* /usr/bin/v2ray-plugin && \
 chmod +x /usr/bin/v2ray-plugin

ENTRYPOINT [ "ssserver", "--log-without-time", "-c", "/etc/shadowsocks-rust/config.json" ]
EOF

docker build -t ssserver-rust-v2ray:latest .


docker rm -f ssserver-rust-v2ray && \
docker run --name ssserver-rust-v2ray \
--restart always \
-p 8389:8388/tcp \
-p 8389:8388/udp \
-v /root/ss/ssconfig.json:/etc/shadowsocks-rust/config.json \
-v "/root/.acme.sh/${host}":"/root/.acme.sh/${host}" \
-dit ssserver-rust-v2ray:latest \
&& docker logs ssserver-rust-v2ray

cat > ssudpconfig.json <<- EOF
{
    "server": "0.0.0.0",
    "server_port": 8388,
    "password": "${password}",
    "timeout": 300,
    "plugin": "v2ray-plugin",
    "plugin_opts": "server;mode=quic;host=${host}",
    "method": "aes-256-gcm"
}
EOF

docker run --name ssserver-rust-quic \
--restart always \
-p 8390:8388/tcp \
-p 8390:8388/udp \
-v /root/ss/ssudpconfig.json:/etc/shadowsocks-rust/config.json \
-v "/root/.acme.sh/${host}":"/root/.acme.sh/${host}" \
-dit ssserver-rust-v2ray:latest \
&& docker logs ssserver-rust-quic