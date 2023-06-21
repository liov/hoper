kubectl apply -f - <<EOF
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
  generation: 1
  name: apisix
  namespace: fino
spec:
  rules:
    - host: gateway.fino
      http:
        paths:
          - backend:
              serviceName: apisix-gateway
              servicePort: 9080
            path: /
EOF