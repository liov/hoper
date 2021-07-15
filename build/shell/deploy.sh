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

CMD ["./${project}"]
EOF

docker_reg=$(docker info |grep reg.*.so$)

tag=${docker_reg}/liov/${project}:v${_version}-$(date "+%Y%m%d%H%M%S")

docker build -t ${tag} .

cat > Dockerfile <<- EOF
apiVersion: v1
kind: Pod
metadata:
  name: ${project}
  labels:
    app: ${project}
spec:
  containers:
    - name: ${project}
      image: ${tag}
EOF

echo "123456" | sudo -S kubectl apply -f ${project}.yaml
sudo kubectl get pod

cd ${current_dir}