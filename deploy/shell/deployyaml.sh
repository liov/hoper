#!/bin/bash

# 设置默认值
confdir=/var/sdk/config
datadir=/var/data

# 使用getopts解析参数
while getopts "f:a:i:c:d:" opt; do
  case $opt in
    f) filepath="$OPTARG" ;;
    a)
       app="$OPTARG"
       group="$OPTARG" ;;
    i) image="$OPTARG" ;;
    c) confdir="$OPTARG" ;;
    d) datadir="$OPTARG" ;;
    \?) echo "Invalid option -$OPTARG" >&2 ;;
  esac
done

# 参数校验
if [ -z "$filepath" ]; then
  echo "必须指定文件路径参数 (-f)" >&2
  exit 1
fi

if [ -z "$app" ]; then
  echo "必须指定应用名称参数 (-a)" >&2
  exit 1
fi

if [ -z "$image" ]; then
  echo "必须指定镜像名称参数 (-i)" >&2
  exit 1
fi

cat <<EOF > $filepath
apiVersion: apps/v1
kind: Deployment
metadata:
 labels:
   app: ${app}
   group: ${group}
 name: ${app}
 namespace: default
spec:
 replicas: 1
 selector:
   matchLabels:
     app: ${app}
 minReadySeconds: 5
 strategy:
   type: RollingUpdate
   rollingUpdate:
     maxSurge: 1
     maxUnavailable: 1
 template:
   metadata:
     labels:
       app: ${app}
       group: ${group}
   spec:
     containers:
       - name: ${app}
         image: ${image}
         resources:
           requests:
             memory: "10Mi"
             cpu: "10m"
           limits:
             memory: "50Mi"
         imagePullPolicy: IfNotPresent
         volumeMounts:
           - mountPath: /app/config
             name: config
           - mountPath: /data
             name: data
     volumes:
       - name: config
         hostPath:
           path: ${confdir}
           type: DirectoryOrCreate
       - name: data
         hostPath:
           path: ${datadir}/${group}
           type: DirectoryOrCreate


EOF