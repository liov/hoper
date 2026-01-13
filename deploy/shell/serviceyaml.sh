#!/bin/bash

while getopts "f:a:p:" opt; do
  case $opt in
    f)
      filepath="$OPTARG"
      ;;
    a)
      app="$OPTARG"
      ;;
    p)
      port="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

cat <<EOF > $filepath
  apiVersion: v1
  kind: Service
  metadata:
    name: ${app}
    namespace: default
    labels:
      app: ${app}
  spec:
    type: ClusterIP
    ports:
      - name: http
        port: 80
        protocol: TCP
        targetPort: ${port}
    selector:
      app: ${app}
EOF