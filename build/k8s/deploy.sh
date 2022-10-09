# docker
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io


---

sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
 sudo apt-get update
 sudo apt-get install docker-ce docker-ce-cli containerd.io docker-compose-plugin


sudo gpasswd -a ${USER} docker
newgrp docker
sudo service docker restart
mkdir /etc/docker
vi /etc/docker/daemon.json
{
    "registry-mirrors": ["https://docker.mirrors.ustc.edu.cn"],
    "insecure-registries":["${ip}"],
}

docker login -u 用户名 -p 密码 ${ip}
# k8s
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
sudo install minikube-linux-amd64 /usr/local/bin/minikube

minikube start --driver=none  --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers --extra-config=apiserver.service-node-port-range=1-39999  --extra-config=kube-proxy.mode=ipvs --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0 --bootstrapper=kubeadm

cp -r /var/lib/minikube/certs/etcd /root/certs/ && chmod 666 /root/certs/etcd/server.key
minikube addons enable dashboard
minikube addons enable logviewer
minikube addons enable efk
minikube addons enable helm-tiller

# prometheus
kubectl create namespace monitoring
helm install kube-prometheus prometheus-community/kube-prometheus-stack -f helm.yaml -n monitoring
# apisix
kubectl create namespace ingress-apisix
# acme
kubectl apply -f tls.yaml
# tools
kubectl create namespace tools

