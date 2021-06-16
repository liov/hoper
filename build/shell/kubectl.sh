# 扩容
kubectl  scale deployment nginx-deployment --replicas=4
# 重启
kubectl get pod pod名称 -n 命名空间名称 -o yaml | kubectl replace --force -f -.

kubectl create configmap my-config-2 --from-file=/etc/resolv.conf
# key的名称是文件名称，value的值是这个文件的内容

kubectl create configmap my-config-3 --from-file=test
# 目录中的文件名为key，文件内容是value

# log
# shellcheck disable=SC2006
kubectl config use-context stage -n namespace&& pod=`kubectl get pods |grep -oE podname[a-zA-Z0-9-]*` && kubectl logs -f $pod

pods=$(kubectl get pods --selector=job-name=pi --output=jsonpath={.items..metadata.name}) && kubectl logs $pods

kubectl describe deployment data-center

pods=$(kubectl get pods --selector=app=${PWD##*/} --output=jsonpath={.items..metadata.name}) && kubectl logs -f $pods

# deploy
../deploy/main -flow all -env dev -name ${PWD##*/} -ns openmng -path . -ver v1.1.0-$(date "+%Y%m%d%H%M%S")

env=stage && git pull && make config && make deploy env=$env tag=v$(date "+%y%m%d%H%M")



# proxy
kubectl proxy --address='0.0.0.0'  --accept-hosts='^*$'