#!/usr/bin/env bash

# shellcheck disable=SC2006
pod=`kubectl get pods |grep -oE customer[a-zA-Z0-9-]*` && kubectl logs -f $pod
