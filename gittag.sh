dir=$(cd $(dirname $0);pwd)
cd $dir/thirdparty/protobuf
go get github.com/hopeio/utils@main
git add .
git commit -m "fix: upgrade dependency"
git commit --amend --date="$(date -d '-10 hours' '+%Y-%m-%d %H:%M:%S')" --no-edit
git tag "v0.0.0"

cd $dir/thirdparty/context
go get github.com/hopeio/utils@main
git add .
git commit -m "fix: upgrade dependency"
git commit --amend --date="$(date -d '-10 hours' '+%Y-%m-%d %H:%M:%S')" --no-edit
git tag "v0.0.0"