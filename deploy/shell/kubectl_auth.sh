#!/bin/bash
server=$1
cluster=$2
kubeconfig="--kubeconfig=/root/.kube/config"

mkdir cert/$cluster
echo $CACRT |base64 -d > cert/$cluster/ca.crt
echo $DEVCRT |base64 -d > cert/$cluster/dev.crt
echo $DEVKEY |base64 -d > cert/$cluster/dev.key


kubectl config set-cluster k8s --server=$server --certificate-authority=cert/$cluster/ca.crt --embed-certs=true kubeconfig
kubectl config set-credentials dev --client-certificate=cert/$cluster/dev.crt --client-key=cert/$cluster/dev.key --embed-certs=true kubeconfig
kubectl config set-context dev --cluster=k8s --user=dev $kubeconfig
kubectl config use-context dev $kubeconfig