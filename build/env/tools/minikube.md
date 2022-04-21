choco install minikube
curl kubectl

minikube ssh
sudo passwd docker


curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
# 特定版本
curl -LO https://dl.k8s.io/release/v1.23.0/bin/linux/amd64/kubectl
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl
-- 挂载data目录不成功，可能是权限问题
docker中部署Kubernetes
minikube start --driver=docker --insecure-registry= --registry-mirror= --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers --mount --mount-string=$HOME:/host --cpus=4 --memory='8192M'
root 直接部署
minikube start --driver=none --registry-mirror= --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers
-- port
--extra-config=apiserver.service-node-port-range=1-10000
-- prometheus-operator
--extra-config=apiserver.authorization-mode=RBAC
--kube-prometheus
--bootstrapper=kubeadm --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0
外网通过代理访问docker中的服务
--url只打印url不自动打开浏览器
//通过代理暴露集群内ip
kubectl proxy --port=8001 --address='0.0.0.0' --accept-hosts='^.*' &
curl http://[k8s-proxy-ip]:8001/api/v1/namespaces/[namespace-name]/services/[service-name]:80/proxy
curl http://[k8s-proxy-ip]:8001/api/v1/namespaces/[namespace-name]/pods/[pod-name]:8080/proxy
& 号将命令放到后台运行
http://localhost:8001/api/v1/namespaces/kubernetes-dashboard/services/http:kubernetes-dashboard:/proxy/#/overview?namespace=default

// 挂载目录
9P Mounts
9P mounts are flexible and work across all hypervisors, but suffers from performance and reliability issues when used with large folders (>600 files). See Driver Mounts as an alternative.

To mount a directory from the host into the guest using the mount subcommand:

minikube mount <source directory>:<target directory>
For example, this would mount your home directory to appear as /host within the minikube VM:

minikube mount $HOME:/host

# helm

$ curl -fsSL -o get_helm.sh https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3
$ chmod 700 get_helm.sh
$ ./get_helm.sh