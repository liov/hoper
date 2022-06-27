需要固定ip,支持k8s再说吧

docker pull apache/doris:build-env-ldb-toolchain-latest

git clone https://github.com/apache/incubator-doris.git

docker run --rm -it -v $PWD/doris:/root/doris  -v $PWD/.m2:/root/.m2 apache/doris:build-env-ldb-toolchain-latest

sh build.sh