FROM node:16-alpine3.15

#修改容器时区
ENV TZ=Asia/Shanghai LANG=C.UTF-8

RUN apk add --update --no-cache \
tzdata && ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

WORKDIR /vhoper

ADD . /vhoper

RUN npm install

RUN npm run build

CMD ["npm","run","start"]