# 删除所有名字中带 “provider”
docker rmi $(docker images | grep "provider" | awk '{print $3}')
