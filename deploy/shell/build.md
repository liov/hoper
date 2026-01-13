docker run --rm  -v $PWD:/work -w /work go build -trimpath main.go

docker run --rm  --privileged=true -u root -v $PWD:/work -w /work node:22-alpine3.16  pnpm run build

kubectl create configmap ${app} --from-file=config.toml,local.toml

kubectl expose deployment ${app} --type=ClusterIP --port=80 --target-port=8080 --name=${app} -n ${namespace}
