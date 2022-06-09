#!/bin/bash

# 具名参数
show_usage="args:[-v, -d] [--version=, --dir=]"

# 版本
_version=1
# 目录
_dir=""

GETOPT_ARGS=`getopt -o l:r: -al version:,dir: -- "$@"`
eval set -- "$GETOPT_ARGS"

#获取参数
while [ -n "$1" ]
    do
        case "$1" in
            -v|--version) _version=$2; shift 2;;
            -d|--dir) _dir=$2; shift 2;;
            --) break ;;
            *) echo $1,$2,$show_usage; break ;;
        esac
done

#对必填项做输入检查，此处假设都为必填项
if [[ -z _version || -z _dir ]]; then
    echo $show_usage
    echo "version: $_version, dir: $_dir"
    exit 0
fi
current_dir=${PWD}

cd ${_dir}

project=${PWD##*/}


GOOS=linux GOARCH=amd64 go build -o ${project}

cat > Dockerfile <<- EOF
FROM frolvlad/alpine-glibc:latest

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /hoper

ADD ./${project} /hoper/${project}
ADD ./config.toml /hoper/config.toml

CMD ["./${project}","-env","test"]
EOF

docker_reg="reg.hoper.xyz"

tag=${docker_reg}/liov/${project}:v${_version}-$(date "+%Y%m%d%H%M%S")

docker build -t ${tag} .

cat > ${project}.yaml <<- EOF
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${project}
  namespace: ${project}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${project}
  minReadySeconds: 5
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 1
  template:
    metadata:
      namespace: deafult
      labels:
        app: ${project}
    spec:
      containers:
        - name: data-center
          image: ${tag}
          imagePullPolicy: Always
          resources:
            # keep request = limit to keep this container in guaranteed class
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 10m
              memory: 50Mi
          volumeMounts:
            - name: logs
              mountPath: /hoper/logs
            - name: config
              mountPath: /hoper
      volumes:
        - name: logs
          hostPath:
            path: /home/crm/dev/logs
        - name: config
          configMap:
            items:
              - key: config
                path: /home/crm/dev/config.toml

---
apiVersion: v1
kind: Service
metadata:
  name: ${project}
  labels:
    app: ${project}
spec:
  selector:
    app: ${project}
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 8090
---
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
  name: ${project}
  namespace: default
spec:
  http:
    - name: ${project}
      match:
        hosts:
          - ${project}.local.org
        paths:
          - /*
      backends:
        - serviceName: ${project}
          servicePort: 80
      plugins:
        - name: cors
          enable: true
EOF

echo "123456" | sudo -S kubectl apply -f ${project}.yaml
sudo kubectl get pod

cd ${current_dir}