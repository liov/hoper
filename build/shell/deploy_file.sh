#!/bin/bash

# 具名参数
show_usage="args:[-v, -f] [--version=, --file=]"

# 版本
_version=1
# 目录
_file=""

GETOPT_ARGS=`getopt -o l:r: -al version:,file: -- "$@"`
eval set -- "$GETOPT_ARGS"

#获取参数
while [ -n "$1" ]
    do
        case "$1" in
            -v|--version) _version=$2; shift 2;;
            -f|--file) _file=$2; shift 2;;
            --) break ;;
            *) echo $1,$2,$show_usage; break ;;
        esac
done

#对必填项做输入检查，此处假设都为必填项
if [[ -z _version || -z _file ]]; then
    echo $show_usage
    echo "version: $_version, file: $_file"
    exit 0
fi


cat > Dockerfile <<- EOF
FROM frolvlad/alpine-glibc:latest

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

ADD ./${_file} /app/${_file}

CMD ["./${_file}","-env","test"]
EOF

docker_reg="reg.hoper.xyz"

tag=${docker_reg}/liov/${_file}:v${_version}-$(date "+%Y%m%d%H%M%S")

docker build -t ${tag} .

cat > ${_file}.yaml <<- EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${_file}
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${_file}
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  template:
    metadata:
      namespace: default
      labels:
        app: ${_file}
    spec:
      containers:
        - name: ${_file}
          image: ${tag}
          imagePullPolicy: IfNotPresent
          resources:
            # keep request = limit to keep this container in guaranteed class
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 100m
              memory: 64Mi
          volumeMounts:
            - name: logs
              mountPath: /app/logs
            - name: static
              mountPath: /app/static
            - name: config
              mountPath: /app
      volumes:
        - name: logs
          hostPath:
            path: /home/crm/dev/logs
        - name: static
          hostPath:
            path: /home/crm/dev/static
        - name: config
          configMap:
            items:
              - key: config
                path: config.toml

---
apiVersion: v1
kind: Service
metadata:
  name: ${_file}
  labels:
    app: ${_file}
spec:
  selector:
    app: ${_file}
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 8090
---
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
  name: ${_file}
  namespace: default
spec:
  http:
    - name: ${_file}
      match:
        hosts:
          - ${_file}.local.org
        paths:
          - /*
      backends:
        - serviceName: ${_file}
          servicePort: 80
      plugins:
        - name: cors
          enable: true
EOF

echo "123456" | sudo -S kubectl apply -f ${_file}.yaml
sudo kubectl get pod
