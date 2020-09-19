tar -chzf zoneinfo.tar.gz /usr/share/zoneinfo/
docker build -t reg.hoper.xyz/liov/attendance:1.6 .
docker tag reg.hoper.xyz/liov/attendance:1.6 reg.hoper.xyz/liov/attendance
docker push  reg.hoper.xyz/liov/attendance:latest
minikube ssh
docker pull reg.hoper.xyz/liov/actix:latest
alias kube=kubectl
kube get pod
kube delete -f attendance.yaml
kube create -f attendance.yaml

minikube start --insecure-registry=reg.hoper.xyz --driver=docker --registry-mirror=https://registry.docker-cn.com --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers
minikube start --insecure-registry=reg.hoper.xyz --driver=docker -- --memory='8192M' --cpus=4

ping registry.aliyuncs.com/google_containers
ping registry.aliyuncs.com
