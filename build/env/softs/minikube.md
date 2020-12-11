choco install minikube
curl kubectl
minikube start
minikube ssh
sudo passwd docker

docker中部署Kubernetes
minikube start --driver=docker  --cpus=4 --memory='8192M' --insecure-registry= --registry-mirror=https://registry.docker-cn.com --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers
外网通过代理访问docker中的服务
--url只打印url不自动打开浏览器
//通过代理暴露集群内ip
kubectl proxy --port=8001 --address='0.0.0.0' --accept-hosts='^.*' &
& 号将命令放到后台运行
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/overview?namespace=default
