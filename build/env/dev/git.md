git verify-pack -v .git/objects/pack/pack-8eaeb...9e.idx | sort -k 3 -n | tail -3

git rev-list --objects --all | grep 35047899fd3b0dd637b0da2086e7a70fe27b1ccb

#git log --pretty=oneline --branches --  

git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *" --prune-empty --tag-name-filter cat -- --all

git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *" --prune-empty $commit-id..HEAD

git rm -rf .git/refs/original/
 
git reflog expire --expire=now --all
 
git fsck --full --unreachable
 
git repack -A -d
 
git gc --aggressive --prune=now
 
git push --force

git filter-repo

git fetch --all && git reset --hard origin/master && git pull

使用 git config --system --unset credential.helper 方法 清除保存好的账号密码
git config --global credential.helper store

git checkout -b xxx= git branch xxx && git checkout xxx
git add xxx.xxx
git commit -m 'xxx'

git remote add

如果要同步你的工作，运行 git fetch origin 命令。 这个命令查找 “origin” 是哪一个服务器（在本例中，它是 git.ourcompany.com），从中抓取本地没有的数据，并且更新本地数据库，移动 origin/master 指针指向新的、更新后的位置。

//rebase
git checkout 次分支
git rebase master //此时次分支最新
//一般我们这样做的目的是为了确保在向远程分支推送时能保持提交历史的整洁——例如向某个其他人维护的项目贡献代码时
git rebase origin master

git checkout master
git merge 次分支

假设你希望将 client 中的修改合并到主分支并发布，但暂时并不想合并 server 中的修改，因为它们还需要经过更全面的测试。 这时，你就可以使用 git rebase 命令的 --onto 选项，选中在 client 分支里但不在 server 分支里的修改（即 C8 和 C9），将它们在 master 分支上重放：

$ git rebase --onto master server client
以上命令的意思是：“取出 client 分支，找出处于 client 分支和 server 分支的共同祖先之后的修改，然后把它们在 master 分支上重放一遍”。 这理解起来有一点复杂，不过效果非常酷。

There is no tracking information for the current branch.
Please specify which branch you want to rebase against.
See git-rebase(1) for details.

    git rebase '<branch>'

If you wish to set tracking information for this branch you can do so with:

    git branch --set-upstream-to=<remote>/<branch> master

rebase最简单的理解应该是将本分支的修改应用在变基的分支上，变基的分支是不变的

//fork同步，添加远程库合并
对fork的代码进行同步更新：

git remote -v #查看当前项目的远程仓库配置
git remote add upstream 原始项目仓库的git地址 # 把原项目的远程仓库添加到fork的代码的远程中
git remote -v # 可以看到原项目的远程仓库已经在配置里了
4.git fetch upstream # 拉取最新的代码
5. git merge upstream/master # mege

git remote add origin xxx
git remote set-url origin xxx
git pull origin master --allow-unrelated-histories
git remote show origin
git push --set-upstream origin master

Git恢复reset --hard丢失的文件
通常最快捷的办法是使用 git reflog 工具。当你 (在一个仓库下) 工作时，Git 会在你每次修改了 HEAD 时悄悄地将改动记录下来。当你提交或修改分支时，reflog 就会更新。git update-ref 命令也可以更新 reflog。
1.先用reflog看看记录的所有HEAD的历史: git reflog

2.然后找到那个SHA，进行恢复: git reset --hard 98abc5a

git fsck --lost-found碰运气


# 修改某次注释
git rebase -i HEAD~2

你想修改哪条注释，就把哪条注释前面的pick换成edit，注意不要动注释内容，只要改前面的东西就好了。
i进入编辑模式，把pick换成edit后，Esc退出编辑模式，:wq保存并退出。

接下来输入：
git commit --amend

修改好注释内容后，输入：
git rebase --continue


# 更改历史提交人信息
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

# windows下必须双引号
git filter-branch -f --env-filter "GIT_COMMITTER_EMAIL=lby.i@qq.com;GIT_AUTHOR_EMAIL=lby.i@qq.com" 1cd75a677457d08c803e40e7d4f317c957cc8562..HEAD

# 代理
git config --global http.proxy 'socks5://127.0.0.1:1080'

git config --global https.proxy 'socks5://127.0.0.1:1080'

git config --global --unset http.proxy

git config --global --unset https.proxy
vi ~/.gitconfig
[http]
proxy = socks5://127.0.0.1:2080
[https]
proxy = socks5://127.0.0.1:2080

# 多个远程仓库
方法 1：每次push、pull时需分开操作
首先，查看本地仓库所关联的远程仓库：（假定最初仅关联了一个远程仓库）

$ git remote -v
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)

然后，用git remote add 添加一个远程仓库，其中name可以任意指定（对应上面的origin部分），比如：

$ git remote add coding.net git@git.coding.net:KeithNull/keithnull.github.io.git

再次查看本地仓库所关联的远程仓库，可以发现成功关联了两个远程仓库：

$ git remote -v
coding.net      git@git.coding.net:KeithNull/keithnull.github.io.git (fetch)
coding.net      git@git.coding.net:KeithNull/keithnull.github.io.git (push)
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)

此后，若需进行push操作，则需要指定目标仓库，git push ，对这两个远程仓库分别操作：

$ git push origin master
$ git push coding.net master

同理，pull操作也需要指定从哪个远程仓库拉取，git pull ，从这两个仓库中选择其一：

$ git pull origin master
$ git pull coding.net master

方法 2：push和pull无需额外操作
在方法 1 中，由于我们添加了多个远程仓库，在push和pull时便面临了仓库的选择问题。诚然如此较为严谨，但是在许多情况下，我们只需要保持远程仓库完全一致，而不需要进行区分，因而这样的区分便显得有些“多余”。

同样地，先查看已有的远程仓库：（假定最初仅关联了一个远程仓库）

$ git remote -v
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)
然后，不额外添加远程仓库，而是给现有的远程仓库添加额外的 URL。使用git remote set-url -add ，给已有的名为name的远程仓库添加一个远程地址，比如：

$ git remote set-url --add origin git@git.coding.net:KeithNull/keithnull.github.io.git
再次查看所关联的远程仓库：

$ git remote -v
origin  git@github.com:keithnull/keithnull.github.io.git (fetch)
origin  git@github.com:keithnull/keithnull.github.io.git (push)
origin  git@git.coding.net:KeithNull/keithnull.github.io.git (push)
可以看到，我们并没有如方法 1 一般增加远程仓库的数目，而是给一个远程仓库赋予了多个地址（或者准确地说，多个用于push的地址）。

因此，这样设置后的push 和pull操作与最初的操作完全一致，不需要进行调整。

删除所有远程标签
git show-ref --tag | awk '{print ":" $2}' | xargs git push origin

删除所有本地标签
git tag -l | xargs git tag -d