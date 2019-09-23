FROM golang:latest
WORKDIR app
# 将当前目录内容复制到位于 /hoper 中的容器中
ADD ./go app/hoper.xyz/hoper
RUN go build . -o hoper
EXPOSE 8000
ENTRYPOINT [ "./hoper" ]