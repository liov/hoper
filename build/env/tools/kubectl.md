# kubectl[https://kubernetes.io/docs/tasks/tools/install-kubectl/]
```bash
linux:
curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/linux/amd64/kubectl
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version

curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
curl -LO "https://dl.k8s.io/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl.sha256"
echo "$(<kubectl.sha256)  kubectl" | sha256sum --check
sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

ubuntu:
sudo apt-get update && sudo apt-get install -y apt-transport-https
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo apt-key add -
echo "deb https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee -a /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

sudo apt-get install -y apt-transport-https ca-certificates curl
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
sudo apt-get install -y kubectl

win:
https://storage.googleapis.com/kubernetes-release/release/v1.16.0/bin/windows/amd64/kubectl.exe
```
## 配置kubectl
```bash
mkdir $HOME/.kube/${cluster}

# 配置一个名为 ${cluster} 的集群，并指定服务地址与根证书
kubectl config set-cluster ${cluster} --server=${cluster-server} --certificate-authority=$HOME/.kube/${cluster}/ca.pem
# 设置一个用户为 ${Username} ，并配置访问的授权文件
kubectl config set-credentials ${Username} --client-certificate=${name}.crt --client-key=${name}.pem --embed-certs=true
# 设置一个名为 ${context} 使用 ${cluster} 集群与 ${Username} 用户的上下文，
kubectl config set-context ${context} --cluster=${cluster} --user=${Username} [--namespace=]
# 启用 ${context} 
kubectl config use-context ${context}

${cluster} 集群名称，可自己定义
${cluster-server} 集群地址，取值如下： 
prod: https://hoper.xyz
test: https:/hoper.xyz
dev: https://hoper.xyz
${Username} 用户名称，自行定义,注意不同context的用户名不能相同，否则会覆盖配置
${name} 邮件回复的授权文件名，可自行修改，文件名与配置一致即可
${context} 上下文，自行定义
```