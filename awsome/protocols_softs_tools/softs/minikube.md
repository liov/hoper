choco install minikube
curl kubectl
minikube start
minikube ssh
sudo passwd docker

docker中部署Kubernetes
minikube start --driver=docker 
外网通过代理访问docker中的服务
--url只打印url不自动打开浏览器
//暴露给外网
minikube kubectl --  proxy --port=8001 --address='192.168.1.212' --accept-hosts='^.*' &