git clone https://github.com/apache/rocketmq-operator
cd rocketmq-operator && make deploy
cd example && kubectl create -f rocketmq_v1alpha1_rocketmq_cluster.yaml