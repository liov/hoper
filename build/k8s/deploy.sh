# docker
sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl gnupg-agent software-properties-common
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo apt-key fingerprint 0EBFCD88
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
sudo apt-get update
sudo apt-get install docker-ce docker-ce-cli containerd.io
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

minikube start --driver=none --registry-mirror= --image-repository=registry.cn-hangzhou.aliyuncs.com/google_containers --extra-config=apiserver.service-node-port-range=1-65536  --extra-config=kubelet.authentication-token-webhook=true --extra-config=kubelet.authorization-mode=Webhook --extra-config=scheduler.bind-address=0.0.0.0 --extra-config=controller-manager.bind-address=0.0.0.0 --bootstrapper=kubeadm

minikube addons enable dashboard
minikube addons enable logviewer
minikube addons enable efk
minikube addons enable helm-tiller

kubectl edit cm kube-proxy -n kube-system
mode 改为 ipvs
kubectl get pod -n kube-system | grep kube-proxy |awk '{system("kubectl delete pod "$1" -n kube-system")}'

# helm

# apisix
helm repo add apisix https://charts.apiseven.com
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
kubectl create ns ingress-apisix
helm install apisix apisix/apisix \
  --set gateway.type=NodePort \
  --set ingress-controller.enabled=true \
  --namespace ingress-apisix \
  --set ingress-controller.config.apisix.serviceNamespace=ingress-apisix
kubectl get service --namespace ingress-apisix

# minikube
kubectl edit StatefulSet apisix-etcd -n ingress-apisix
replicas: 1
kubectl delete PersistentVolumeClaim data-apisix-etcd-1 -n ingress-apisix
kubectl delete PersistentVolumeClaim data-apisix-etcd-2 -n ingress-apisix

helm repo add apisix https://charts.apiseven.com
helm repo update
helm install apisix-dashboard apisix/apisix-dashboard --namespace ingress-apisix

vim apisix-dashboard.yaml - |
apiVersion: apisix.apache.org/v2beta3
kind: ApisixRoute
metadata:
  name: apisix-dashboard
  namespace: ingress-apisix
spec:
  http:
    - name: apisix-dashboard
      match:
        hosts:
          - apisix.hoper.xyz
        paths:
          - /*
      backends:
        - serviceName: apisix-dashboard
          servicePort: 80
          resolveGranularity: service
kubectl apply -f apisix-dashboard.yaml
vim /home/ubuntu/deploy/hoper/nginx/nginx.conf
server {
        listen       80;
        server_name  localhost *.hoper.xyz;
        location / {
                    proxy_pass  http://127.0.0.1:30687;
                    proxy_set_header Host $http_host;
                    proxy_set_header X-Real-IP $remote_addr;
                    proxy_set_header X-Real-PORT $remote_port;
                    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
               }
        location /api/live/ws {
            proxy_pass http://127.0.0.1:30687;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection "Upgrade";
            proxy_set_header Host  $http_host;
        }
        error_page  404 403 500 502 503 504  /error;
    }
/usr/local/openresty/nginx/sbin/nginx -c /home/ubuntu/deploy/hoper/nginx/nginx.conf -s reload

kubectl edit cm apisix -n ingress-apisix
apisix:
 enable_control: true
  control:
    ip: "127.0.0.1"
    port: 9090
  plugin_attr:
        prometheus:
          export_uri: /metrics
          export_addr:
            ip: 0.0.0.0
            port: 9091

kubectl edit deployment apisix -n ingress-apisix
- containerPort: 9090
  name: control
  protocol: TCP
- containerPort: 9091
  name: prometheus
  protocol: TCP

# prometheus
# grafana
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install kube-prometheus prometheus-community/kube-prometheus-stack -n monitoring
helm pull prometheus-community/kube-prometheus-stack
helm install kube-prometheus kube-prometheus-stack-34.10.0.tgz -n monitoring
helm upgrade kube-prometheus prometheus-community/kube-prometheus-stack -n monitoring
-- uninstall
helm uninstall kube-prometheus -n monitoring
kubectl delete crd alertmanagerconfigs.monitoring.coreos.com
kubectl delete crd alertmanagers.monitoring.coreos.com
kubectl delete crd podmonitors.monitoring.coreos.com
kubectl delete crd probes.monitoring.coreos.com
kubectl delete crd prometheuses.monitoring.coreos.com
kubectl delete crd prometheusrules.monitoring.coreos.com
kubectl delete crd servicemonitors.monitoring.coreos.com
kubectl delete crd thanosrulers.monitoring.coreos.com

docker pull bitnami/kube-state-metrics:2.4.1
docker tag bitnami/kube-state-metrics:2.4.1 k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.4.1
docker pull registry.aliyuncs.com/prometheus-adapter/prometheus-adapter:v0.9.1
docker pull willdockerhub/prometheus-adapter:v0.9.1
docker tag docker.io/willdockerhub/prometheus-adapter:v0.9.1 k8s.gcr.io/prometheus-adapter/prometheus-adapter:v0.9.1
docker pull bitnami/kube-state-metrics:2.4.2
docker tag docker.io/bitnami/kube-state-metrics:2.4.2 k8s.gcr.io/kube-state-metrics/kube-state-metrics:v2.4.2
docker pull dyrnq/kube-webhook-certgen:v1.1.1
docker tag dyrnq/kube-webhook-certgen:v1.1.1 k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.1.1
docker pull rancher/curlimages-curl:7.73.0

kubectl get secret kube-prometheus-grafana -n monitoring -o yaml
echo 'cHJvbS1vcGVyYXRvcg==' | base64 --decode
prom-operator
# 仓库安装
git clone https://github.com/prometheus-operator/kube-prometheus
kubectl apply --server-side -f manifests/setup
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
kubectl apply -f manifests/

kubectl apply --server-side -f manifests/setup -f manifests
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup

cat > prometheus-additional.yaml <<- EOF
 - job_name: "apisix-prometheus"
   scrape_interval: 10s
   metrics_path: "/apisix/prometheus/metrics"
   static_configs:
   - targets: ["apisix-prometheus.ingress-apisix:9091"]
EOF

kubectl create secret generic additional-scrape-configs --from-file=prometheus-additional.yaml --dry-run=client -oyaml > additional-scrape-configs.yaml

kubectl apply -f additional-scrape-configs.yaml -n monitoring
cat >  prometheus.yaml <<- EOF
apiVersion: monitoring.coreos.com/v1
kind: Prometheus
metadata:
  name: prometheus
  labels:
    prometheus: prometheus
spec:
  replicas: 2
  serviceAccountName: prometheus
  serviceMonitorSelector:
    matchLabels:
      team: frontend
  additionalScrapeConfigs:
    name: additional-scrape-configs
    key: prometheus-additional.yaml
EOF

import Apache APISIX 11719
https://github.com/apache/apisix-ingress-controller/blob/master/docs/assets/other/json/apisix-ingress-controller-grafana.json
# 第一个图
Total Requests
sum(apisix_http_requests_total{instance=~"$instance"})

