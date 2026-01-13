#!/bin/bash

while getopts "f:a:c:r:" opt; do
  case $opt in
    f) filepath="$OPTARG" ;;
    a) app="$OPTARG" ;;
    c) cmd="$OPTARG" ;;
    r) register="${OPTARG}/" ;;
    \?) echo "Invalid option -$OPTARG" >&2 ;;
  esac
done

cat <<EOF > $filepath
FROM ${register}jybl/timezone AS tz

FROM ${register}frolvlad/alpine-glibc

#修改容器时区
ARG TZ=Asia/Shanghai
ENV TZ=\${TZ} LANG=C.UTF-8
COPY --from=tz /usr/share/zoneinfo/\$TZ /usr/share/zoneinfo/\$TZ
RUN echo \$TZ > /etc/timezone && ln -sf /usr/share/zoneinfo/\$TZ /etc/localtime

WORKDIR /app

ADD ./${app} /app

CMD [$cmd]
EOF