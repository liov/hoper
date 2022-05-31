#!/bin/bash

echo $CA > ca.crt
echo $CACRT > dev.crt
echo $CAKEY > dev.key

server=https://hoper.xyz:6443
kubectl config set-cluster k8s --server=${server} --certificate-authority=ca.crt --embed-certs=true --kubeconfig=/root/dev.conf
kubectl config set-credentials dev --client-certificate=dev.crt --client-key=dev.key --embed-certs=true --kubeconfig=/root/dev.conf
kubectl config set-context dev --cluster=k8s --user=dev --kubeconfig=/root/dev.conf
kubectl config use-context dev --kubeconfig=/root/dev.conf