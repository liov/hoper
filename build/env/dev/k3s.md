wget https://github.com/rancher/k3s/releases/download/v0.5.0/k3s
chmod 777 k3s
export http_proxy=http://ip:port
export https_proxy=http://ip:port
source /etc/profile
curl -s https://zhangguanzhang.github.io/bash/pull.sh | bash -s --  k8s.gcr.io/pause:3.1
docker tag k8s.gcr.io/pause:3.1 liovjyb/pause:3.1
docker login
docker push liovjyb/pause:3.1
sudo ./k3s server --pause-image=liovjyb/pause:3.1
https://github.com/rancher/k3s/issues/396

docker pull mirrorgooglecontainers/xxx:vx.y.z

docker tag mirrorgooglecontainers/xxx:vx.y.z k8s.gcr.io/xxx:vx.y.z