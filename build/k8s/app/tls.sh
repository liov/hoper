openssl rand -hex 16

# openssl
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=hoper.xyz/O=hoper.xyz"


# acme
docker run --rm  -it  \
  -v $PWD/acme:/acme.sh  \
  --net=host \
  neilpang/acme.sh --issue --dns -d hoper.xyz -d *.hoper.xyz  --yes-I-know-dns-manual-mode-enough-go-ahead-please

docker run --rm  -it  \
  -v $PWD/acme:/acme.sh  \
  --net=host  \
  neilpang/acme.sh --renew  -d hoper.xyz -d *.hoper.xyz \
    --yes-I-know-dns-manual-mode-enough-go-ahead-please

kubectl delete secret hoper-xyz -n ingress-apisix
kubectl create secret tls hoper-xyz -n ingress-apisix --cert=fullchain.cer --key=hoper.xyz.key

kubectl apply -f - <<EOF
apiVersion: apisix.apache.org/v2beta3
kind: ApisixTls
metadata:
  name: hoper-xyz
  namespace: ingress-apisix
spec:
  hosts:
    - hoper.xyz
    - "*.hoper.xyz"
  secret:
    name: hoper-xyz
    namespace: ingress-apisix
EOF

# cert-manager
kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: ca-issuer
  namespace: default
spec:
  ca:
    secretName: ca-key-pair
EOF

kubectl apply -f - <<EOF
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: hoper-xyz
  namespace: default
spec:
  # Secret names are always required.
  secretName: hoper-xyz-tls

  # secretTemplate is optional. If set, these annotations and labels will be
  # copied to the Secret named example-com-tls. These labels and annotations will
  # be re-reconciled if the Certificate's secretTemplate changes. secretTemplate
  # is also enforced, so relevant label and annotation changes on the Secret by a
  # third party will be overwriten by cert-manager to match the secretTemplate.
  secretTemplate:
    annotations:
      my-secret-annotation-1: "foo"
      my-secret-annotation-2: "bar"
    labels:
      my-secret-label: foo

  duration: 2160h # 90d
  renewBefore: 360h # 15d
  subject:
    organizations:
      - jetstack
  # The use of the common name field has been deprecated since 2000 and is
  # discouraged from being used.
  commonName: example.com
  isCA: false
  privateKey:
    algorithm: RSA
    encoding: PKCS1
    size: 2048
  usages:
    - server auth
    - client auth
  # At least one of a DNS Name, URI, or IP address is required.
  dnsNames:
    - hoper.xyz
  uris:
    - spiffe://cluster.local/ns/sandbox/sa/example
  ipAddresses:
    - 192.168.0.5
  # Issuer references are always required.
  issuerRef:
    name: ca-issuer
    # We can reference ClusterIssuers by changing the kind here.
    # The default value is Issuer (i.e. a locally namespaced Issuer)
    kind: Issuer
    # This is optional since cert-manager will default to this value however
    # if you are using an external issuer, change this to that issuer group.
    group: cert-manager.io
EOF