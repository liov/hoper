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
  git push origin $new_version
}

cd $dir/thirdparty/initialize
upgrade "utils"
cd $dir/thirdparty/protobuf
upgrade "utils"
cd $dir/thirdparty/context
upgrade "utils"

cd $dir/thirdparty/pick
upgrade "context"