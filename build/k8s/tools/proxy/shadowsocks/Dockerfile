FROM ghcr.io/shadowsocks/ssserver-rust:latest AS ssserver

USER root

RUN cd /tmp && \
 TAG=$(wget -qO- https://api.github.com/repos/shadowsocks/v2ray-plugin/releases/latest | grep tag_name | cut -d '"' -f4) && \
 wget https://github.com/shadowsocks/v2ray-plugin/releases/download/$TAG/v2ray-plugin-linux-amd64-$TAG.tar.gz && \
 tar -xf *.gz && \
 rm *.gz && \
 mv v2ray* /usr/bin/v2ray-plugin && \
 chmod +x /usr/bin/v2ray-plugin

USER nobody

ENTRYPOINT [ "ssserver", "--log-without-time", "-c", "/etc/shadowsocks-rust/config.json" ]

FROM ghcr.io/shadowsocks/sslocal-rust:latest AS sslocal

USER root

RUN cd /tmp && \
 TAG=$(wget -qO- https://api.github.com/repos/shadowsocks/v2ray-plugin/releases/latest | grep tag_name | cut -d '"' -f4) && \
 wget https://github.com/shadowsocks/v2ray-plugin/releases/download/$TAG/v2ray-plugin-linux-amd64-$TAG.tar.gz && \
 tar -xf *.gz && \
 rm *.gz && \
 mv v2ray* /usr/bin/v2ray-plugin && \
 chmod +x /usr/bin/v2ray-plugin

USER nobody

ENTRYPOINT [ "sslocal", "--log-without-time", "-c", "/etc/shadowsocks-rust/config.json" ]