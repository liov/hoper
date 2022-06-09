# istio[https://github.com/istio/istio/releases/tag/1.3.0]
```bash
linux:
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=1.3.0 sh -
cd istio-1.3.0
export PATH=$PWD/bin:$PATH
```
# istio[https://istio.io/docs/setup/install/kubernetes/]
```bash
linux:
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl apply -f $i; done
win ps:
Get-ChildItem e:/istio-1.3.0/install/kubernetes/helm/istio-init/files/ | ForEach-Object -Process{
    if ($_.Name -match "crd*yaml"])
    {
        kubectl apply -f $_;
    }

}

kubectl apply -f install/kubernetes/istio-demo.yaml
kubectl label namespace default istio-injection=enabled
kubectl apply -f samples/bookinfo/platform/kube/bookinfo.yaml
```
## 卸载istio
```bash
kubectl delete -f install/kubernetes/istio-demo.yaml
for i in install/kubernetes/helm/istio-init/files/crd*yaml; do kubectl delete -f $i; done
```

## 部署istio
```bash
当您使用时部署应用程序时kubectl apply，如果Istio边车注入器 在标有istio-injection=enabled以下标记的名称空间中启动，它们将自动将Envoy容器注入您的应用程序窗格：
kubectl label namespace <namespace> istio-injection=enabled
kubectl create -n <namespace> -f <your-app-spec>.yaml
在没有istio-injection标签的命名空间中，您可以使用 istioctl kube-inject 在部署它们之前在应用程序窗格中手动注入Envoy容器：
istioctl kube-inject -f <your-app-spec>.yaml | kubectl apply -f -
```

# 2020-06-15
```bash
cd istio-1.6.2
export PATH=$PWD/bin:$PATH
istioctl manifest apply --set profile=demo
```
## 卸载
```bash
istioctl manifest generate --set profile=demo | kubectl delete -f -
```

# Helm 部署(已弃用)
```bash
curl -L https://git.io/getLatestIstio | sh -
sudo apt-get install -y jq
ISTIO_VERSION=$(curl -L -s https://api.github.com/repos/istio/istio/releases/latest | jq -r .tag_name)
cd istio-${ISTIO_VERSION}
cp bin/istioctl /usr/local/bin

kubectl create -f install/kubernetes/helm/helm-service-account.yaml
helm init --service-account tiller

kubectl apply -f install/kubernetes/helm/istio/templates/crds.yaml
helm install install/kubernetes/helm/istio --name istio --namespace istio-system \
  --set ingress.enabled=true \
  --set gateways.enabled=true \
  --set galley.enabled=true \
  --set sidecarInjectorWebhook.enabled=true \
  --set mixer.enabled=true \
  --set prometheus.enabled=true \
  --set grafana.enabled=true \
  --set servicegraph.enabled=true \
  --set tracing.enabled=true \
  --set kiali.enabled=false
```
# 2022-06-09
curl -L https://istio.io/downloadIstio | sh -
cd istio-1.14.0
export PATH=$PWD/bin:$PATH

istioctl install --set profile=default -y （我觉得不用部署ingress）
给命名空间添加标签，指示 Istio 在部署应用的时候，自动注入 Envoy 边车代理：
kubectl label namespace default istio-injection=enabled
设置入站端口：
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export SECURE_INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}'
设置入站 IP：
export INGRESS_HOST=$(minikube ip)
设置环境变量 GATEWAY_URL:
export GATEWAY_URL=$INGRESS_HOST:$INGRESS_PORT

# 查看仪表板
kubectl apply -f samples/addons
kubectl rollout status deployment/kiali -n istio-system
istioctl dashboard kiali