helm repo add apisix https://charts.apiseven.com
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
kubectl create ns ingress-apisix
```bash
helm install apisix apisix/apisix \
  --set gateway.type=NodePort \
  --set ingress-controller.enabled=true \
  --namespace ingress-apisix \
  --set ingress-controller.config.apisix.serviceNamespace=ingress-apisix
```
`kubectl get service --namespace ingress-apisix`
`helm upgrade apisix apisix/apisix --install -n ingress-apisix`

`helm repo add apisix https://charts.apiseven.com`
`helm repo update`
`helm install apisix-dashboard apisix/apisix-dashboard --namespace ingress-apisix`
```bash
vim apisix-dashboard.yaml - |
apiVersion: apisix.apache.org/v2
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
```


# new
https://artifacthub.io/packages/helm/apisix/apisix/0.9.3

dashboard false 没有挂载目录，tls的目录

`kubectl create secret generic etcd-ssl --from-file=/root/certs/etcd/ca.crt -n ingress-apisix`
`kubectl create secret generic ssl --from-file=/root/deploy/app/hoper/acme/hoper.xyz/ca.cer -n ingress-apisix`

```bash
helm install apisix apisix/apisix \
  -f helm.yaml \
  --namespace ingress-apisix
```
```bash
helm upgrade apisix apisix/apisix \
  -f helm.yaml \
  --namespace ingress-apisix \
  --install
```

`kubectl edit cm apisix -n ingress-apisix`
```yaml
node_listen:       # APISIX listening port
      - port: 9081
        enable_http2: true
```
`kubectl edit svc apisix-gateway -n ingress-apisix`
```yaml
name: prometheus
port: 9091
protocol: TCP
targetPort: 9091
```
如果设置了serviceMonitor.namespace
`kubectl edit ServiceMonitor apisix -n monitoring`
```yaml
namespaceSelector:
    matchNames:
    - ingress-apisix
```