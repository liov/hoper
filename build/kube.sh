#!/bin/bash
cluster=dev
server=https://hoper.xyz:6666
Username=liov
name=liov
context=dev
# 配置一个名为 ${cluster} 的集群，并指定服务地址与根证书
kubectl config set-cluster ${cluster} --server=${server} --certificate-authority=$HOME/.kube/${cluster}/${cluster}-ca.pem
# 设置一个用户为 ${Username} ，并配置访问的授权文件
kubectl config set-credentials ${Username} --client-certificate=$1${name}.crt --client-key=$1${name}.pem --embed-certs=true
# 设置一个名为 ${context} 使用 ${cluster} 集群与 ${Username} 用户的上下文，
kubectl config set-context ${context} --cluster=${cluster} --user=${Username} --namespace=$2
# 启用 ${context} 
kubectl config use-context ${context}