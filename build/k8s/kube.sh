#!/bin/bash
cluster=$1
if [[ $1 == dev ]];then
	server=https://k8s.dev:6443
elif [[ $1 == test ]]; then
	server=https://k8s.test:6443
elif [[ $1 == prod ]]; then
	server=https://k8s.prod
fi
Username=liov
name=liov
# 配置一个名为 ${cluster} 的集群，并指定服务地址与根证书
kubectl config set-cluster ${cluster} --server=${server} --certificate-authority=./${cluster}-ca.pem
# 设置一个用户为 ${Username} ，并配置访问的授权文件
kubectl config set-credentials ${Username} --client-certificate=./${cluster}-${name}.crt --client-key=./${cluster}-${name}.pem --embed-certs=true
# 设置一个名为 ${context} 使用 ${cluster} 集群与 ${Username} 用户的上下文，
kubectl config set-context ${cluster} --cluster=${cluster} --user=${Username}
# 启用 ${context} 
kubectl config use-context ${cluster}