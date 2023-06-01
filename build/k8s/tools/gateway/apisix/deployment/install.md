helm repo add apisix https://charts.apiseven.com
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
kubectl create ns ingress-apisix




# new
https://artifacthub.io/packages/helm/apisix/apisix/0.9.3

dashboard false 没有挂载目录，tls的目录
cp -r /var/lib/minikube/certs/etcd /root/certs/ && chmod 666 /root/certs/etcd/server.key || k8s.runAsUser=0 || initContainer.command - chown -R nobody:nobody /certs/etcd
kubectl create secret generic etcd-ssl --from-file=/root/certs/etcd/ -n ingress-apisix
kubectl create secret generic ssl --from-file=/root/certs/acme/hoper.xyz/ca.cer -n ingress-apisix

```bash
helm install apisix apisix/apisix -f helm.yaml -n ingress-apisix
kubectl get pod -n ingress-apisix
```
# upgrade
```bash
helm upgrade apisix apisix/apisix -f helm.yaml -n ingress-apisix --install
```

# 如果设置了serviceMonitor
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
