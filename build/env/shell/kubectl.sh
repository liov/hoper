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
kubectl config use-context stage -n namespace && kubectl logs -f $(kubectl get pods |grep -oE podname[a-zA-Z0-9-]+)

kubectl logs $(kubectl get pods --selector=job-name=pi --output=jsonpath={.items..metadata.name})

kubectl describe deployment data-center

kubectl logs -f $(kubectl get pods --selector=app=${PWD##*/} --output=jsonpath={.items..metadata.name})

# deploy
../deploy/main -flow all -env dev -name ${PWD##*/} -ns ${USER} -path . -ver v1.1.0-$(date "+%Y%m%d%H%M%S")

env=stage && git pull && make config && make deploy env=$env tag=v$(date "+%y%m%d%H%M")



# proxy
kubectl proxy --address='0.0.0.0'  --accept-hosts='^.*$'
# logs
Examples:
  # Return snapshot logs from pod nginx with only one container
  kubectl logs nginx

  # Return snapshot logs from pod nginx with multi containers
  kubectl logs nginx --all-containers=true

  # Return snapshot logs from all containers in pods defined by label app=nginx
  kubectl logs -lapp=nginx --all-containers=true

  # Return snapshot of previous terminated ruby container logs from pod web-1
  kubectl logs -p -c ruby web-1

  # Begin streaming the logs of the ruby container in pod web-1
  kubectl logs -f -c ruby web-1

  # Begin streaming the logs from all containers in pods defined by label app=nginx
  kubectl logs -f -lapp=nginx --all-containers=true

  # Display only the most recent 20 lines of output in pod nginx
  kubectl logs --tail=20 nginx

  # Show all logs from pod nginx written in the last hour
  kubectl logs --since=1h nginx

  # Return snapshot logs from first container of a job named hello
  kubectl logs job/hello

  # Return snapshot logs from container nginx-1 of a deployment named nginx
  kubectl logs deployment/nginx -c nginx-1

Options:
      --all-containers=false: Get all containers logs in the pod(s).
  -c, --container='': Print the logs of this container
  -f, --follow=false: Specify if the logs should be streamed.
      --ignore-errors=false: If watching / following pod logs, allow for any errors that occur to be non-fatal
      --limit-bytes=0: Maximum bytes of logs to return. Defaults to no limit.
      --max-log-requests=5: Specify maximum number of concurrent logs to follow when using by a selector. Defaults to 5.
      --pod-running-timeout=20s: The length of time (like 5s, 2m, or 3h, higher than zero) to wait until at least one pod is running
  -p, --previous=false: If true, print the logs for the previous instance of the container in a pod if it exists.
  -l, --selector='': Selector (label query) to filter on.
      --since=0s: Only return logs newer than a relative duration like 5s, 2m, or 3h. Defaults to all logs. Only one of since-time / since may be used.
      --since-time='': Only return logs after a specific date (RFC3339). Defaults to all logs. Only one of since-time / since may be used.
      --tail=-1: Lines of recent log file to display. Defaults to -1 with no selector, showing all log lines otherwise 10, if a selector is provided.


# 将pod进行缩容操作 让其为0 即等同于停止操作

kubectl scale --replicas=0 deployment/<your-deployment>

kubectl get pod -n kube-system | grep kube-proxy |awk '{system("kubectl delete pod "$1" -n kube-system")}'

# windows
kubectl proxy --port=8001 --address=0.0.0.0 --accept-hosts=^.* --kubeconfig=