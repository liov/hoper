kubectl create namespace polardbx-operator-system
helm install --namespace polardbx-operator-system polardbx-operator https://github.com/ApsaraDB/galaxykube/releases/download/v1.2.1/polardbx-operator-1.2.1.tgz

helm repo add polardbx https://polardbx-charts.oss-cn-beijing.aliyuncs.com
helm install --namespace polardbx-operator-system polardbx-operator polardbx/polardbx-operator

echo "apiVersion: polardbx.aliyun.com/v1
kind: PolarDBXCluster
metadata:
  name: quick-start
  annotations:
    polardbx/topology-mode-guide: quick-start" | kubectl apply -f -

kubectl delete polardbxcluster quick-start
helm uninstall --namespace polardbx-operator-system polardbx-operator