openssl passwd [-apr1] <password>

kubectl delete cm nginx && kubectl create configmap nginx --from-file=nginx

kubectl delete cm nginx && kubectl create configmap nginx --from-file=nginx && pods=$(kubectl get pods --selector=app=openresty --output=jsonpath={.items..metadata.name}) && kubectl exec $pods -- openresty -c /usr/local/openresty/nginx/conf/nginx.conf -s reload