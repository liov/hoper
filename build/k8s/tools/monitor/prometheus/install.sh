# prometheus
# grafana
## helm kube-prometheus-stack
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm install kube-prometheus prometheus-community/kube-prometheus-stack -f helm.yaml -n monitoring
helm pull prometheus-community/kube-prometheus-stack
helm install kube-prometheus kube-prometheus-stack-34.10.0.tgz   -f helm.yaml -n monitoring
helm upgrade kube-prometheus prometheus-community/kube-prometheus-stack  -f helm.yaml -n monitoring
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
echo 'cHJvbS1vcGVyYXRvcg==' | base64 -d
prom-operator


## 仓库安装 prometheus-operator 别试，坑
git clone https://github.com/prometheus-operator/kube-prometheus
kubectl apply --server-side -f manifests/setup
until kubectl get servicemonitors --all-namespaces ; do date; sleep 1; echo ""; done
kubectl apply -f manifests/

kubectl apply --server-side -f manifests/setup -f manifests
kubectl delete --ignore-not-found=true -f manifests/ -f manifests/setup

# 外部访问
kubectl edit cm kube-prometheus-grafana -n monitoring
[security]
allow_embedding = true
[auth.anonymous]
enabled = true
kubectl delete pod kube-prometheus-grafana-5b8cbc5d5b-t44zc -n monitoring

helm upgrade kube-prometheus prometheus-community/kube-prometheus-stack -f helm.yaml -n monitoring