#!/bin/bash

dir=${1%/*}
echo "dir: $dir"
file=${1##*/}
echo "file: $file"

echo "chmod +x $1"
chmod +x $1

echo "docker stop $file"
docker stop $file
echo "docker rm $file"
docker rm $file
docker rmi $(docker images | grep $file | awk '{print $3}') -f

# Dockerfile
cat > $dir/Dockerfile <<- EOF
FROM frolvlad/alpine-glibc:latest

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

ADD ./${file} /app/${file}

CMD ["./${file}"]
EOF


image=$file:$(date "+%Y%m%d%H%M%S")
echo "docker build -t $image $dir"
docker build -t $image $dir
echo "docker run --name $pod -d -v /data/$file:/app/data $image"
docker run --name $file -d -v /data/$file:/app/data -v $dir/config:/app/config $image

echo "docker logs $file"
docker logs $file