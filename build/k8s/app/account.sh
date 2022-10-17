#!/bin/bash
dir=$1
cluster=$2
if [[ $2 == tx ]];then
	server=https://hoper.xyz:6443
	echo $CA |base64 -d > ca.crt
    echo $CACRT |base64 -d > dev.crt
    echo $CAKEY |base64 -d > dev.key
elif [[ $2 == tot ]]; then
	server=https://192.168.1.212:6443
	cd $dir/$cluster
fi

kubectl config set-cluster k8s --server=${server} --certificate-authority=ca.crt --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-credentials dev --client-certificate=dev.crt --client-key=dev.key --embed-certs=true --kubeconfig=/root/.kube/config
kubectl config set-context dev --cluster=k8s --user=dev --kubeconfig=/root/.kube/config
kubectl config use-context dev --kubeconfig=/root/.kube/config