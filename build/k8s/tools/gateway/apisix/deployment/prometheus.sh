kubectl edit cm apisix -n ingress-apisix
```yaml
apisix:
 enable_control: true
  control:
    ip: "127.0.0.1"
    port: 9090
  plugin_attr:
    prometheus:
      export_uri: /apisix/prometheus/metrics
      export_addr:
        ip: 0.0.0.0
        port: 9091
```
kubectl edit deployment apisix -n ingress-apisix
```yaml
- containerPort: 9090
  name: control
  protocol: TCP
- containerPort: 9091
  name: prometheus
  protocol: TCP
```
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

# APISIX
import Apache APISIX 11719
# 第一个图
Total Requests
sum(apisix_http_requests_total{instance=~"$instance"})
# APISIX Controller
https://github.com/apache/apisix-ingress-controller/blob/master/docs/assets/other/json/apisix-ingress-controller-grafana.json