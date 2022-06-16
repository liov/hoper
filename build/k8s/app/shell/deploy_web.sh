#!/bin/bash

# 具名参数
show_usage="args:[-v, -p] [--version=, --project=]"

# 版本
_version=1
# 目录
_project=""

GETOPT_ARGS=`getopt -o l:r: -al version:,project: -- "$@"`
eval set -- "$GETOPT_ARGS"

#获取参数
while [ -n "$1" ]
    do
        case "$1" in
            -v|--version) _version=$2; shift 2;;
            -p|--project) _project=$2; shift 2;;
            --) break ;;
            *) echo $1,$2,$show_usage; break ;;
        esac
done

#对必填项做输入检查，此处假设都为必填项
if [[ -z _version || -z _project ]]; then
    echo $show_usage
    echo "version: $_version, project: $_project"
    exit 0
fi


tag=${_project}-front-end

cat > ${tag}.yaml <<- EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-configmap
data:
  nginx_conf: |-
    #user  nobody;
    worker_processes  1;
    events {
        worker_connections  1024;
    }
    http {
        include       mime.types;
        default_type  application/octet-stream;
        sendfile        on;
        keepalive_timeout  65;
        gzip  on;

        server {
            listen       80;
            server_name  localhost;
            location /static/ {
                root  /hoper/static/;
            }
            location / {
                root   /hoper/dist/;
                index  index.html index.htm;
            }
            error_page   500 502 503 504  /50x.html;
            location = /50x.html {
                root   html;
            }
        }
    }
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: ${tag}
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: ${tag}
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
        app: ${tag}
    spec:
      containers:
        - name: ${tag}
          image: nginx:latest
          imagePullPolicy: IfNotPresent
          resources:
            # keep request = limit to keep this container in guaranteed class
            limits:
              cpu: 500m
              memory: 512Mi
            requests:
              cpu: 10m
              memory: 50Mi
          volumeMounts:
            - name: dist
              mountPath: /hoper/dist
            - name: static
              mountPath: /hoper/static
            - name: nginx
              subPath: nginx.conf
              mountPath: /etc/nginx/nginx.conf
      volumes:
        - name: dist
          hostPath:
            path: /home/dev/dev/dist
        - name: static
          hostPath:
            path: /home/dev/dev/static
        - name: nginx
          configMap:
            name: nginx-configmap
            items:
              - key: nginx_conf
                path: nginx.conf

---
apiVersion: v1
kind: Service
metadata:
  name: ${tag}
  labels:
    app: ${tag}
spec:
  selector:
    app: ${tag}
  ports:
    - name: http
      protocol: TCP
      port: 80
      targetPort: 80
---
apiVersion: apisix.apache.org/v2alpha1
kind: ApisixRoute
metadata:
  name: ${tag}
  namespace: default
spec:
  http:
    - name: ${tag}
      match:
        hosts:
          - ${_project}.local.org
        paths:
          - /*
      backends:
        - serviceName: ${tag}
          servicePort: 80
      plugins:
        - name: cors
          enable: true
EOF

echo "123456" | sudo -S kubectl apply -f ${tag}.yaml
sudo kubectl get pod
