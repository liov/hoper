#!/bin/bash
buildDir="build"
if [ ! -d $buildDir ]; then
    mkdir $buildDir
fi

# 使用getopts处理参数
while getopts ":t:s:r:" opt; do
  case ${opt} in
    t )
      target=$OPTARG
      ;;
    s )
      source_path=$OPTARG
      ;;
    r )
      register="$OPTARG/"
      ;;
    \? )
      echo "无效参数: -$OPTARG" 1>&2
      exit 1
      ;;
    : )
      echo "参数 -$OPTARG 需要一个值" 1>&2
      exit 1
      ;;
  esac
done


# build
echo GOOS=linux go build -trimpath -o build/${target} $source_path
GOOS=linux go build -trimpath -o  build/${target} $source_path
cmd="\"./${target}\",\"-c\",\"./config/${target}.toml\""
dockerfilepath=build/Dockerfile
rundir=$(realpath $(dirname "${BASH_SOURCE[0]}"))
echo $rundir
source ${rundir}/dockerfile.sh -f $dockerfilepath -a $target -c $cmd -r $register

#docker run --rm -v $GOPATH:/go -v $PWD:/work -w /work -e GOPROXY=$GOPROXY $GOIMAGE go build  -trimpath -o /work/build/$output /work/$target
image=${register}jybl/$target
source ${rundir}/deployyaml.sh build/${target}.yaml $target $image
source ${rundir}/serviceyaml.sh -f build/${target}_service.yaml -a $target -p 9000
echo "docker build -t $image -f $dockerfilepath $buildDir; docker push $image"

if command -v wsl >/dev/null 2>&1; then
   wsl bash -c "cd /mnt/$PWD; pwd; docker build -t $image -f $dockerfilepath $buildDir; docker push $image"
else
   docker build -t $image -f $dockerfilepath $buildDir; docker push $image
fi

