dir=$(cd $(dirname $0);pwd)

function upgrade(){
  local last_tag=$(git describe --tags --abbrev=0)
  local patch=$(echo $last_tag | cut -d'.' -f3)
  local new_version="${last_tag%.*}.$((patch + 1))"
  echo $new_version
  git add .
  git commit -m "chore: upgrade dependency"
  git commit --amend --date="$(date -d '-10 hours' '+%Y-%m-%d %H:%M:%S')" --no-edit
  git tag "$new_version"
}

# 获取参数值
param=$1

# 根据参数执行不同的逻辑
case $param in
    "")
        echo "Parameter is empty."
        # 在这里执行空参数的逻辑
        cd $dir/thirdparty/context
        go get github.com/hopeio/utils@main
        upgrade
        cd $dir/thirdparty/initialize
        go get github.com/hopeio/utils@main
        upgrade
        cd $dir/thirdparty/protobuf
        go get github.com/hopeio/utils@main
        upgrade
        cd $dir/thirdparty/deploy/plugin/go
        go get github.com/hopeio/utils@main
        cd $dir/thirdparty/deploy
        upgrade
        ;;
    pc)
        echo "Parameter is pc."
        # 在这里执行pc参数的逻辑
        cd $dir/thirdparty/pick
        go get github.com/hopeio/context@main
        upgrade
        cd $dir/thirdparty/cherry
        go get github.com/hopeio/protobuf@main
        go get github.com/hopeio/context@main
        upgrade
        ;;
    co)
        echo "Parameter is co."
        # 在这里执行co参数的逻辑
        cd $dir/thirdparty/collection
        go get github.com/hopeio/cherry@main
        go get github.com/hopeio/pick@main
        upgrade
        ;;
    *)
        echo "Invalid parameter: $param"
        exit 1
        ;;
esac