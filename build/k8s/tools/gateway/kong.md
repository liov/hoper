kong.conf
```properties
database = postgres
pg_host = postgre.tools
pg_port = 5432
pg_user = web
pg_password = 123456
pg_database = openmng-gw
admin_listen = 0.0.0.0:8001, 0.0.0.0:8444 ssl
plugins = bundled,session,request-inspector,session-go
lua_package_path = /usr/local/?.lua;/usr/local/?/init.lua;

pluginserver_names = go

pluginserver_go_socket = /usr/local/kong/go_pluginserver.sock
pluginserver_go_start_cmd = /usr/local/bin/go-pluginserver -kong-prefix /usr/local/kong/ -plugins-directory /usr/local/kong/go-plugins
pluginserver_go_query_cmd = /usr/local/bin/go-pluginserver -dump-all-plugins -plugins-directory /usr/local/kong/go-plugins
nginx_user = root
```
```bash
#!/usr/bin/env bash

cd go-plugins/session-go
go build  -trimpath github.com/Kong/go-pluginserver
go build  -trimpath -o session-go.so -buildmode plugin session_validator.go check_path.go config.go exchange_token.go redis.go
cd ../../

cat > Dockerfile <<- EOF
FROM kong:2.5-centos

USER root

RUN rm /etc/localtime
RUN ln -s /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

ADD ./deploy/${1}/kong/kong.conf /etc/kong/

EOF

for dir in plugins/*
do
if [ -d ${dir} ]
then
    echo "ADD ./${dir} /usr/local/share/lua/5.1/kong/${dir}"  >> Dockerfile
fi
done

echo "ADD ./go-plugins/session-go/go-pluginserver /usr/local/bin/go-pluginserver"  >> Dockerfile
echo "ADD ./go-plugins/session-go/session-go.so /usr/local/kong/go-plugins/"  >> Dockerfile

docker build . -t $2
```
