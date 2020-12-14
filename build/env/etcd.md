etcdctl --endpoints=https://127.0.0.1:2379 endpoint status --cacert=/var/lib/minikube/certs/etcd/ca.crt --cert=/var/lib/minikube/certs/etcd/server.crt --key=/var/lib/minikube/certs/etcd/server.key

etcd --advertise-client-urls=https://192.168.49.2:2379 \
    --cert-file=/var/lib/minikube/certs/etcd/server.crt \
    --data-dir=/var/lib/minikube/etcd \
    --initial-advertise-peer-urls=https://192.168.49.2:2380 \
    --initial-cluster=minikube=https://192.168.49.2:2380 \
    --key-file=/var/lib/minikube/certs/etcd/server.key \
    --listen-client-urls=https://127.0.0.1:2379,https://192.168.49.2:2379 \
    --listen-metrics-urls=http://127.0.0.1:2381 \
    --listen-peer-urls=https://192.168.49.2:2380 \
    --name=minikube \
    --peer-cert-file=/var/lib/minikube/certs/etcd/peer.crt \
    --peer-key-file=/var/lib/minikube/certs/etcd/peer.key \
    --peer-trusted-ca-file=/var/lib/minikube/certs/etcd/ca.crt \
    --proxy-refresh-interval=70000 \
    --snapshot-count=10000 \
    --trusted-ca-file=/var/lib/minikube/certs/etcd/ca.crt \
    
etcdctl --endpoints=http://192.168.1.212:2379 endpoint status 