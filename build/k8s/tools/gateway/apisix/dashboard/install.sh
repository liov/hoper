helm repo add apisix https://charts.apiseven.com
helm repo update
helm install apisix-dashboard apisix/apisix-dashboard -f dhelm.yaml -n ingress-apisix

kubectl edit cm apisix-dashboard -n ingress-apisix

```yaml
etcd:
mtls:
  key_file: /etcd-ssl/server.key
  cert_file: /etcd-ssl/server.crt
  ca_file: /etcd-ssl/ca.crt
```
kubectl edit deployment apisix-dashboard -n ingress-apisix
volumeMounts:
 - mountPath: /etcd-ssl
   name: etcd-ssl

volumes:
- name: etcd-ssl
  secret:
    defaultMode: 420
    secretName: etcd-ssl

name: etcd-certs



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