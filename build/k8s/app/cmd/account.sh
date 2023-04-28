#!/bin/bash
cluster=$1
if [[ $cluster == tx ]];then
	server=https://hoper.xyz:6443
	mkdir ../certs/$cluster
	echo $CACRT |base64 -d > ../certs/$cluster/ca.crt
    echo $DEVCRT |base64 -d > ../certs/$cluster/dev.crt
    echo $DEVKEY |base64 -d > ../certs/$cluster/dev.key
elif [[ $1 == tot ]]; then
	server=https://192.168.1.212:6443
	cd certs/$cluster
fi

cd ../certs
server=https://hoper.xyz:6443
kubectl config set-cluster k8s --server=$server --certificate-authority=$cluster/ca.crt --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-credentials dev --client-certificate=$cluster/dev.crt --client-key=$cluster/dev.key --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-context dev --cluster=k8s --user=dev --kubeconfig=/root/.kube/config
kubectl config use-context dev --kubeconfig=/root/.kube/config