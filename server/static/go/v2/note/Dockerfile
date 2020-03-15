FROM golang:latest

#修改容器时区
RUN apt update && apt install -y tzdata
ENV TZ Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN && dpkg-reconfigure -f noninteractive tzdata

WORKDIR app
# 将当前目录内容复制到位于 /hoper 中的容器中
ADD ./go app/hoper.xyz/hoper
RUN go build . -o hoper
EXPOSE 8000
ENTRYPOINT [ "./hoper" ]