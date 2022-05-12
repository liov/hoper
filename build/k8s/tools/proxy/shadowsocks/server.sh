{
    "server": "0.0.0.0",
    "server_port": 8388,
    "password": "password",
    "timeout": 300,
    "method": "aes-256-gcm"
}

docker run --name ssserver-rust \
--restart always \
-p 6366:8388/tcp \
-p 6366:8388/udp \
-v /root/ssconfig.json:/etc/shadowsocks-rust/config.json \
-dit ghcr.io/shadowsocks/ssserver-rust:latest
