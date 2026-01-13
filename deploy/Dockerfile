#FROM golang:1.22-alpine3.16 AS builder
#
#ENV GOPROXY https://goproxy.io,https://goproxy.cn,direct
#WORKDIR /build
# ADD . /build
#RUN go build -trimpath -o /build/deploy
#
FROM jybl/timezone AS tz
FROM bitnami/kubectl AS kubectl

FROM docker:20.10.19-cli-alpine3.16

#修改容器时区
ARG TZ=Asia/Shanghai
ENV TZ=${TZ} LANG=C.UTF-8
COPY --from=tz /usr/share/zoneinfo/$TZ /usr/share/zoneinfo/$TZ
RUN echo $TZ > /etc/timezone && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime


COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /bin/

ADD ./tpl /tpl
ADD ./shell /shell

