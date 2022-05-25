kubectl create configmap nginx --from-file=nginx

kubectl delete cm nginx && kubectl create configmap nginx --from-file=nginx && kubectl exec openresty-8644644d76-68n2r -- openresty -s reload