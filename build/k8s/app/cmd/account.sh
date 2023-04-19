#!/bin/bash
cluster=$1
if [[ $1 == tx ]];then
	server=https://hoper.xyz:6443
	echo $CACRT |base64 -d > ../certs/$1/ca.crt
    echo $DEVCRT |base64 -d > ../certs/$1/dev.crt
    echo $DEVKEY |base64 -d > ../certs/$1/dev.key
elif [[ $1 == tot ]]; then
	server=https://192.168.1.212:6443
	cd certs/$cluster
fi
server=https://hoper.xyz:6443
kubectl config set-cluster k8s --server=${server} --certificate-authority=certs/$1/ca.crt --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-credentials dev --client-certificate=certs/$1/dev.crt --client-key=certs/$1/dev.key --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-context dev --cluster=k8s --user=dev --kubeconfig=/root/.kube/config
kubectl config use-context dev --kubeconfig=/root/.kube/config