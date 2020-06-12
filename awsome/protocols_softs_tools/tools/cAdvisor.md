Heapster：在k8s集群中获取metrics和事件数据，写入InfluxDB，heapster收集的数据比cadvisor多，却全，而且存储在influxdb的也少。

Heapster将每个Node上的cAdvisor的数据进行汇总，然后导到InfluxDB。

Heapster的前提是使用cAdvisor采集每个node上主机和容器资源的使用情况，
再将所有node上的数据进行聚合。

这样不仅可以看到Kubernetes集群的资源情况，
还可以分别查看每个node/namespace及每个node/namespace下pod的资源情况。
可以从cluster，node，pod的各个层面提供详细的资源使用情况。
InfluxDB：时序数据库，提供数据的存储，存储在指定的目录下。

Grafana：提供了WEB控制台，自定义查询指标，从InfluxDB查询数据并展示。

cAdvisor+Prometheus+Grafana
访问http://localhost:8080/metrics，可以拿到cAdvisor暴露给 Prometheus的数据

