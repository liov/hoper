# Install Ingress APISIX on Minikube

This document explains how to install Ingress APISIX on [Minikube](https://minikube.sigs.k8s.io/).

## Prerequisites

* Install [Minikube](https://minikube.sigs.k8s.io/docs/start/).
* Install [Helm](https://helm.sh/).
* Clone [Apache APISIX Charts](https://github.com/apache/apisix-helm-chart).
* Clone [apisix-ingress-controller](https://github.com/apache/apisix-ingress-controller).
* Make sure your target namespace exists, kubectl operations thorough this document will be executed in namespace `ingress-apisix`.

## Install APISIX

[Apache APISIX](http://apisix.apache.org/) as the proxy plane of apisix-ingress-controller, should be deployed in advance.

```shell
cd /path/to/apisix-helm-chart
helm repo add bitnami https://charts.bitnami.com/bitnami
helm dependency update ./chart/apisix
helm install apisix ./chart/apisix \
  --set allow.ipList="{0.0.0.0/0}" \
  --namespace ingress-apisix
kubectl get service --namespace ingress-apisix
```

Two Service resources were created, one is `apisix-gateway`, which processes the real traffic; another is `apisix-admin`, which acts as the control plane to process all the configuration changes.

One thing should be concerned that the `allow.ipList` field should be customized according to the Pod CIRD configuration of Minikube, so that the apisix-ingress-controller instances can access the APISIX instances (resources pushing).

## Install apisix-ingress-controller

You can also install apisix-ingress-controller by Helm Charts, it's recommended to install it in the same namespace with Apache APISIX.

```shell
cd /path/to/apisix-ingress-controller
# install apisix-ingress-controller
helm install apisix-ingress-controller ./charts/apisix-ingress-controller \
  --set image.tag=dev \
  --set config.apisix.baseURL=http://apisix-admin:9180/apisix/admin \
  --set config.apisix.adminKey=edd1c9f034335f136f87ad84b625c8f1 \
  --namespace ingress-apisix
```

Change the `image.tag` to the apisix-ingress-controller version that you desire. You have to wait for while until the correspdoning pods are running.

Now try to create some [resources](../CRD-specification.md) to verify the running of Ingress APISIX. As a minimalist example, see [proxy-the-httpbin-service](../samples/proxy-the-httpbin-service.md) to learn how to apply resources to drive the apisix-ingress-controller.