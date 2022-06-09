# 开发
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-kubernetes/master/all-in-one/jaeger-all-in-one-template.yml
# 存储
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-kubernetes/master/production-elasticsearch/configmap.yml
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-kubernetes/master/production-elasticsearch/elasticsearch.yml
# 生产
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-kubernetes/master/jaeger-production-template.yml
# Service Dependencies
kubectl run jaeger-spark-dependencies --schedule="55 23 * * *" --env="STORAGE=elasticsearch" --env="ES_NODES=elasticsearch:9200" --env="ES_USERNAME=changeme" --env="ES_PASSWORD=changeme" --restart=Never --image=jaegertracing/spark-dependencies

kubectl create namespace observability # <1>
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/crds/jaegertracing.io_jaegers_crd.yaml # <2>
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/service_account.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/role_binding.yaml
kubectl create -n observability -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/operator.yaml

kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role.yaml
kubectl create -f https://raw.githubusercontent.com/jaegertracing/jaeger-operator/master/deploy/cluster_role_binding.yaml