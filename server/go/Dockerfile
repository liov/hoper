FROM frolvlad/alpine-glibc:latest

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /hoper

ADD ./app /hoper/app

CMD ["./app"]