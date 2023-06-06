
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


kubectl apply -f apisix-dashboard.yaml