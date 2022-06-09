FROM frolvlad/alpine-glibc AS frpc

WORKDIR /app

ADD ./frpc /app

ADD frpc.ini /app

CMD ["./frpc","-c","./frpc.ini"]

FROM frolvlad/alpine-glibc AS frps

WORKDIR /app

ADD ./frps /app

ADD frps.ini /app

CMD ["./frps","-c","./frps.ini"]