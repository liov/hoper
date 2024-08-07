dir=$(cd $(dirname $0);pwd)

function upgrade(){
  local last_tag=$(git describe --tags --abbrev=0)
  local patch=$(echo $last_tag | cut -d'.' -f3)
  local new_version="${last_tag%.*}.$((patch + 1))"
  echo $new_version
  local name="$1"
  go get github.com/hopeio/$name@main
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
          upgrade "utils"
          cd $dir/thirdparty/initialize
          upgrade "utils"
          cd $dir/thirdparty/protobuf
          upgrade "utils"
        ;;
    pick)
        echo "Parameter is pick."
        # 在这里执行pick参数的逻辑
        cd $dir/thirdparty/pick
        upgrade "context"
        ;;
    cherry)
        echo "Parameter is cherry."
        # 在这里执行cherry参数的逻辑

        upgrade "context"
        ;;
    *)
        echo "Invalid parameter: $param"
        exit 1
        ;;
esac