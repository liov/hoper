真心被整无语了，自己执行
https://github.com/alibaba/nacos/blob/master/distribution/conf/nacos-mysql.sql
https://github.com/alibaba/nacos/blob/develop/distribution/conf/schema.sql

git clone https://github.com/nacos-group/nacos-k8s.git

kubectl create -f deploy/nfs/rbac.yaml

# Set the subject of the RBAC objects to the current namespace where the provisioner is being deployed
$ NS=$(kubectl config get-contexts|grep -e "^\*" |awk '{print $5}')
$ NAMESPACE=${NS:-default}
$ sed -i'' "s/namespace:.*/namespace: $NAMESPACE/g" ./deploy/nfs/rbac.yaml


$ NAMESPACE=tools
$ sed -i'' "s/namespace:.*/namespace: $NAMESPACE/g" ./deploy/nfs/rbac.yaml


kubectl create -f deploy/nfs/deployment.yaml
kubectl create -f deploy/nfs/class.yaml
kubectl get pod -l app=nfs-client-provisioner

Modify deploy/nacos/nacos-pvc-nfs.yaml
kubectl create -f nacos-k8s/deploy/nacos/nacos-pvc-nfs.yaml

kubectl create -f nacos-k8s/deploy/nacos/nacos-pvc-nfs.yaml


真坑啊，mysql8 还要自己放驱动
https://repo1.maven.org/maven2/mysql/mysql-connector-java/8.0.29/mysql-connector-java-8.0.29.jar
放plugins/mysql目录下

kubectl scale sts nacos --replicas=1

https://github.com/nacos-group/nacos-template/blob/master/nacos-grafana.json

这内存需求量，简直了， 还没放东西就要爆炸，一个这玩意2G起步，无语
单机集群内存需求差距太大了吧