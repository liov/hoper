# 扩容
kubectl  scale deployment nginx-deployment --replicas=4

kubectl get pod pod名称 -n 命名空间名称 -o yaml | kubectl replace --force -f -.

kubectl create configmap my-config-2 --from-file=/etc/resolv.conf
# key的名称是文件名称，value的值是这个文件的内容

kubectl create configmap my-config-3 --from-file=test
# 目录中的文件名为key，文件内容是value