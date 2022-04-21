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

minikube start  --driver=docker --memory='8192M' --cpus=4 --base-image="kicbase/stable:v0.0.15-snapshot4" --insecure-registry=reg.hoper.xyz --registry-mirror= --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers

ping registry.aliyuncs.com/google_containers
ping registry.aliyuncs.com

docker pull jettech/kube-webhook-certgen:v1.2.2
docker pull siriuszg/nginx-ingress-controller:v0.41.2
docker pull jettech/kube-webhook-certgen:v1.3.0
docker tag jettech/kube-webhook-certgen:v1.3.0 registry.cn-hangzhou.aliyuncs.com/google_containers/kube-webhook-certgen:v1.3.0

