git verify-pack -v .git/objects/pack/pack-8eaeb...9e.idx | sort -k 3 -n | tail -3

git rev-list --objects --all | grep 35047899fd3b0dd637b0da2086e7a70fe27b1ccb

#git log --pretty=oneline --branches --  

git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *" --prune-empty --tag-name-filter cat -- --all

git filter-branch --force --index-filter "git rm --cached --ignore-unmatch *" --prune-empty $commit-id..HEAD

rm -rf .git/refs/original/
 
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