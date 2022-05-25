方法一：
Dockerfile里加上这段:

RUN apk update && apk add tzdata
RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone
方法二：
Dockerfile里加上这段:

RUN apk update && apk add tzdata
RUN cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

方法三：
Dockerfile里加上这段:

RUN echo -n 'VFppZjIAAAAAAAAAAAAAAAAAAAAAAAADAAAAAwAAAAAAAAAdAAAAAwAAAAyAAAAAoJeigKF5BPDIWV6AyQn5cMnTvQDLBYrwy3xAANI7PvDTi3uA1EKt8NVFIgDWTL/w1zy/ANgGZnDZHfKA2UF88B66UiAfaZuQIH6EoCFJfZAiZ6EgIylfkCRHgyAlEnwQJidlICbyXhAoB0cgKNJAEAIBAgECAQIBAgECAQIBAgECAQIBAgECAQIBAgECAABx1wAAAAB+kAEEAABwgAAITE1UAENEVABDU1QAAAAAAAAAVFppZjIAAAAAAAAAAAAAAAAAAAAAAAADAAAAAwAAAAAAAAAdAAAAAwAAAAz/////fjZDKf////+gl6KA/////6F5BPD/////yFlegP/////JCflw/////8nTvQD/////ywWK8P/////LfEAA/////9I7PvD/////04t7gP/////UQq3w/////9VFIgD/////1ky/8P/////XPL8A/////9gGZnD/////2R3ygP/////ZQXzwAAAAAB66UiAAAAAAH2mbkAAAAAAgfoSgAAAAACFJfZAAAAAAImehIAAAAAAjKV+QAAAAACRHgyAAAAAAJRJ8EAAAAAAmJ2UgAAAAACbyXhAAAAAAKAdHIAAAAAAo0kAQAgECAQIBAgECAQIBAgECAQIBAgECAQIBAgECAQIAAHHXAAAAAH6QAQQAAHCAAAhMTVQAQ0RUAENTVAAAAAAAAAAKQ1NULTgK'|base64 -d > /etc/localtime && echo -n 'Asia/Shanghai' > /etc/timezone

方法四：
如果你的镜像已经生成了，那么在启动容器时，可以使用挂载宿主机时区文件的方式，配置镜像时区。当然，镜像的时间也是随着宿主机时间改变的。所以此种方法首先要保证宿主机时间是正确的。

启动docker容器时加上下面这段：

-v /etc/localtime:/etc/localtime -v /etc/timezone:/etc/timezone
