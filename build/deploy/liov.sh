#!/bin/bash

if [ $# == 0 ]
then
   dir=$PWD
   name=${PWD##*/}
   tmp=${PWD##*app/}
   pod=${tmp/\//-}
   if [$name == $tmp]
   then
    pod=$pod-$pod
    fi
fi
if [ $# == 1 ]
then
  dir=~/app/$1
  name=$1
  pod=$1
fi
if [ $# == 2 ]
then
   dir=~/app/$1/$2
   name=$2
   pod=$1-$2
fi

echo "chmod +x $dir/$name"
chmod +x $dir/$name

echo "docker stop $pod"
docker stop $pod
echo "docker rm $pod"
docker rm $pod
docker rmi $(docker images | grep $pod | awk '{print $3}') -f

cat > $dir/Dockerfile <<- EOF
FROM frolvlad/alpine-glibc:latest

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /app

ADD ./${name} /app/${name}

CMD ["./${name}"]
EOF
image=$pod:$(date "+%Y%m%d%H%M%S")
echo "docker build -t $image $dir"
docker build -t $image $dir
echo "docker run --name $pod -d -v /data/$name:/app/data $image"
docker run --name $pod -d -v /data/$name:/app/data -v $dir/config.toml:/app/config.toml -v $dir/local.toml:/app/local.toml -v $dir/config:/app/config $image
echo "docker logs $pod"
docker logs $pod