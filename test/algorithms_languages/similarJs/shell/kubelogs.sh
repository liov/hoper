#!/usr/bin/env bash

# shellcheck disable=SC2006
kubectl config use-context stage -n namespace&& pod=`kubectl get pods |grep -oE podname[a-zA-Z0-9-]*` && kubectl logs -f $pod
