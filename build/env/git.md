[TOC]

git.md
# 一般用法
```bash
git checkout -b xxx= git branch xxx && git checkout xxx
git add xxx.xxx
git commit -m 'xxx'
```

`git verify-pack -v .git/objects/pack/pack-8eaeb...9e.idx | sort -k 3 -n | tail -3`

`git rev-list --objects --all | grep 35047899fd3b0dd637b0da2086e7a70fe27b1ccb`

## 修改历史提交,删除*.jpg
git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *.jpg" --prune-empty --tag-name-filter cat -- --all

git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *.jpg" --prune-empty $commit-id..HEAD

git rm -rf .git/refs/original/

git reflog expire --expire=now --all

git fsck --full --unreachable

git repack -A -d

git gc --aggressive --prune=now

git push --force

git filter-repo

git fetch --all && git reset --hard origin/master && git pull

## 修改某次注释
git rebase -i HEAD~2

你想修改哪条注释，就把哪条注释前面的pick换成edit，注意不要动注释内容，只要改前面的东西就好了。
i进入编辑模式，把pick换成edit后，Esc退出编辑模式，:wq保存并退出。

接下来输入：
git commit --amend

修改好注释内容后，输入：
git rebase --continue


## 更改历史提交人信息
git filter-branch -f --env-filter '
if [ "$GIT_COMMITTER_NAME" = "oldname" ];
then
GIT_COMMITTER_NAME="newname";
GIT_COMMITTER_EMAIL="newaddr";
GIT_AUTHOR_NAME="newname";
GIT_AUTHOR_EMAIL="newaddr";
fi

if [ "$GIT_AUTHOR_NAME" = "oldname" ];
then
GIT_COMMITTER_NAME="newname";
GIT_COMMITTER_EMAIL="newaddr";
GIT_AUTHOR_NAME="newname";
GIT_AUTHOR_EMAIL="newaddr";
fi
' -- --all

### windows下必须双引号
git filter-branch -f --env-filter 'GIT_COMMITTER_NAME="贾一饼";GIT_AUTHOR_NAME="贾一饼";GIT_COMMITTER_EMAIL=lby.i@qq.com;GIT_AUTHOR_EMAIL=lby.i@qq.com' 1cd75a677457d08c803e40e7d4f317c957cc8562..HEAD


如果要同步你的工作，运行 git fetch origin 命令。 这个命令查找 “origin” 是哪一个服务器（在本例中，它是 git.ourcompany.com），从中抓取本地没有的数据，并且更新本地数据库，移动 origin/master 指针指向新的、更新后的位置。

# git rebase
git checkout new
git rebase master //此时new最新
//一般我们这样做的目的是为了确保在向远程分支推送时能保持提交历史的整洁——例如向某个其他人维护的项目贡献代码时
git rebase origin master

git checkout master
git merge new

假设你希望将 client 中的修改合并到主分支并发布，但暂时并不想合并 server 中的修改，因为它们还需要经过更全面的测试。 这时，你就可以使用 git rebase 命令的 --onto 选项，选中在 client 分支里但不在 server 分支里的修改（即 C8 和 C9），将它们在 master 分支上重放：

`git rebase --onto master server client`
以上命令的意思是：“取出 client 分支，找出处于 client 分支和 server 分支的共同祖先之后的修改，然后把它们在 master 分支上重放一遍”。 这理解起来有一点复杂，不过效果非常酷。

There is no tracking information for the current branch.
Please specify which branch you want to rebase against.
See git-rebase(1) for details.

    git rebase '<branch>'

If you wish to set tracking information for this branch you can do so with:

    git branch --set-upstream-to=<remote>/<branch> master

rebase最简单的理解应该是将本分支的修改应用在变基的分支上，变基的分支是不变的

# Git恢复reset --hard丢失的文件
通常最快捷的办法是使用 git reflog 工具。当你 (在一个仓库下) 工作时，Git 会在你每次修改了 HEAD 时悄悄地将改动记录下来。当你提交或修改分支时，reflog 就会更新。git update-ref 命令也可以更新 reflog。
1.先用reflog看看记录的所有HEAD的历史: git reflog

2.然后找到那个SHA，进行恢复: git reset --hard 98abc5a

git fsck --lost-found碰运气

# git config
## 保存账号密码
`git config --system --unset credential.helper` 方法 清除保存好的账号密码
`git config --global credential.helper store` 保存账号密码
## 代理
git config --global http.proxy 'socks5://127.0.0.1:1080'

git config --global https.proxy 'socks5://127.0.0.1:1080'd

git config --global --unset http.proxy

git config --global --unset https.proxy
vi ~/.gitconfig
[http]
proxy = socks5://127.0.0.1:2080
[https]
proxy = socks5://127.0.0.1:2080

# git设置远程仓库地址
git remote set-url origin git@github.com:keithnull/keithnull.github.io.git

## 方法 1：每次push、pull时需分开操作
首先，查看本地仓库所关联的远程仓库：（假定最初仅关联了一个远程仓库）

`git remote -v`
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)

然后，用`git remote add` 添加一个远程仓库，其中name可以任意指定（对应上面的origin部分），比如：

`git remote add coding.net git@git.coding.net:KeithNull/keithnull.github.io.git`

再次查看本地仓库所关联的远程仓库，可以发现成功关联了两个远程仓库：

`git remote -v`
coding.net      git@git.coding.net:KeithNull/keithnull.github.io.git (fetch)
coding.net      git@git.coding.net:KeithNull/keithnull.github.io.git (push)
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)

此后，若需进行push操作，则需要指定目标仓库，git push ，对这两个远程仓库分别操作：

`git push origin master`
`git push coding.net master`

同理，pull操作也需要指定从哪个远程仓库拉取，git pull ，从这两个仓库中选择其一：

`git pull origin master`
`git pull coding.net master`

##  2：push和pull无需额外操作
在方法 1 中，由于我们添加了多个远程仓库，在push和pull时便面临了仓库的选择问题。诚然如此较为严谨，但是在许多情况下，我们只需要保持远程仓库完全一致，而不需要进行区分，因而这样的区分便显得有些“多余”。

同样地，先查看已有的远程仓库：（假定最初仅关联了一个远程仓库）

# git远程仓库
`git remote -v`
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)
然后，不额外添加远程仓库，而是给现有的远程仓库添加额外的 URL。使用git remote set-url -add ，给已有的名为name的远程仓库添加一个远程地址，比如：

`git remote set-url --add origin git@git.coding.net:KeithNull/keithnull.github.io.git`
再次查看所关联的远程仓库：

`git remote -v`
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)
origin  git@git.coding.net:KeithNull/keithnull.github.io.git (push)
可以看到，我们并没有如方法 1 一般增加远程仓库的数目，而是给一个远程仓库赋予了多个地址（或者准确地说，多个用于push的地址）。

因此，这样设置后的push 和pull操作与最初的操作完全一致，不需要进行调整。


## fork同步，添加远程库合并
对fork的代码进行同步更新：
`git remote -v` #查看当前项目的远程仓库配置
`git remote add upstream` 原始项目仓库的git地址 # 把原项目的远程仓库添加到fork的代码的远程中
`git remote -v` # 可以看到原项目的远程仓库已经在配置里了
`git fetch upstream` # 拉取最新的代码
`git merge upstream/master` # merge

git remote add origin xxx
git remote set-url origin xxx
git pull origin master --allow-unrelated-histories
git remote show origin
git push --set-upstream origin master

# git tag
## 删除所有远程标签
`git show-ref --tag | awk '{print ":" $2}' | xargs git push origin`

## 删除所有本地标签
`git tag -l | xargs git tag -d`

## 强制更新
branch=`git rev-parse --abbrev-ref HEAD` && git checkout dev && git merge $branch && git push && git checkout $branch


# git日志
`git log --pretty=oneline --branches`

## git统计某个人代码提交次数
`git log --author="zhangsan" --since='2022-06-10' --until='2022-12-27'  --pretty='%aN' | sort | uniq -c | sort -k1 -n -r`
## git统计某个人代码提交数量
`git log --author=zhangsan --since='2022-01-01 00:00:00' --until='2023-12-31 23:59:59' --pretty=tformat: --numstat | awk '{ add += $1; subs += $2; loc += $1 - $2 } END { printf "added lines: %s, removed lines: %s,total lines: %s\n", add, subs, loc }'`
–author=xxx 查询某一个用户的提交记录
–pretty=tformat: 控制显示的记录格式
–numstat 对增加和删除的行数进行统计 第一列显示的是增加的行数 第二列显示的是删除的行数
–since 需要统计的开始时间
–until 需要统计的结束时间

## git统计所有人代码提交次数
`git log --pretty='%aN' | sort | uniq -c | sort -k1 -n -r`

## git统计所有人代码提交数量
`git log --format='%an' | sort -u | while read name; do echo -en "$name\t"; git log --author="$name" --since=2022-02-07 --until=2023-01-30 --pretty=tformat: --numstat | awk '{ add += $1; subs += $2; loc += $1 - $2 } END { printf "added lines: %s, removed lines: %s, total lines: %s\n", add, subs, loc }' -; done`

## git amend 改提交时间
git commit --amend --date="2024-04-07 19:00:00"
### 固定减去10小时
git commit --amend --date="$(date -d '-10 hours' '+%Y-%m-%d %H:%M:%S')"
### 改提交人
git commit --amend --author="贾一饼 <xxx@.xxx>"
### 修改提交
从HEAD版本开始往过去数3个版本
$ git rebase -i HEAD~3

git rebase -i [commitid]
-i（--interactive）：弹出交互式的界面进行编辑合并
[commitid]：要合并多个版本之前的版本号，注意：[commitid] 本身不参与合并
#### 修改提交
编辑提交前pick->e
git commit --amend xxx
git rebase --continue
#### 合并分支
合并指定版本号（不包含此版本）
$ git rebase -i [commitid]
编辑要合并的提交前pick->s

##### 直接合并?
git commit --fixup=[commitid]
git rebase --autosquash -i


## git config git设置
git config --global user.name "xxx"  --global可选
git config --global user.email "xxx"
git config --global credential.helper store

# git submodule
1. git submodule deinit [<submodule-path>]
2. git submodule add -b <branch> <repository> [<submodule-path>]
3. git submodule update --remote

1.3即可
git submodule foreach -q --recursive 'git checkout $(git config -f $toplevel/.gitmodules submodule.$name.branch || echo master)'