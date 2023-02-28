
docker stop ssserver-rust-v2ray && docker rm ssserver-rust-v2ray
docker run --name ssserver-rust-v2ray --restart always -p 66:8388/tcp -p 66:8388/udp -v /root/ss/ssconfig.json:/etc/shadowsocks-rust/config.json -v /root/ss/cert:/cert -dit ssserver-rust-v2ray:latest && docker logs ssserver-rust-v2ray

docker restart ssserver-rust-v2ray

OVPN_DATA="/root/openvpn"
docker run -v $OVPN_DATA:/etc/openvpn -d -p 1194:1194/udp --cap-add=NET_ADMIN kylemanna/openvpn