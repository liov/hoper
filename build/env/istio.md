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