docker buildx build -t shadowsocks/ssserver-rust:latest -t shadowsocks/ssserver-rust:v1.11.1 --target ssserver .
docker buildx build -t shadowsocks/sslocal-rust:latest -t shadowsocks/sslocal-rust:v1.11.1 --target sslocal .

docker rmi -f ssserver-rust-v2ray && docker build -t ssserver-rust-v2ray:latest .
docker rmi -f sslocal-rust-v2ray && docker build -t sslocal-rust-v2ray:latest .