# 删除所有名字中带 “provider”
docker rmi $(docker images | grep "provider" | awk '{print $3}')
# 查看容器ip
docker inspect d7f29df68dd4 | grep IPAddress
# 删除所有容器
docker rm -f `docker ps -a -q`
# 自动重启
docker run --restart=always
docker update --restart=always <CONTAINER ID>
# root用户
docker run --user="root"
# 覆盖entrypoint
docker run --entrypoint /bin/bash
# 执行后删除
docker run --rm
# dind 特权模式
docker run --user="root" --privileged

docker run --rm -v /mnt/d/SDK/gopath:/go -v $PWD:/work -w /work/server/go/mod golang go run /work/server/go/mod/tools/install.go