作者：Cat Chen
链接：https://www.zhihu.com/question/379545619/answer/1081624808
来源：知乎
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。

作者反思的核心观点是：
万能药是不存在的。
他十年前推荐的 Git Flow 更适用于需要多版本并存的软件，
而如今 web 是潮流且 web 是持续交付的，可能 GitHub Flow 更合适。

我觉得我可以补充一下我的观点。

不仅仅 web 是持续交付的，连 iOS 和 Android 应用本质上也是持续交付的，除非开发者有意按大版本割裂为多个不同的应用。
尽管一个应用可能从 v1.1 升级到 v1.2 后会有人拒绝升级，但本质上 v1.1 和 v1.2 不是多版本并存的。
如果发现了新的 bug 要修复，v1.1 是不能独立升级到 v1.1.1 修复这个 bug 然后避开 v1.2 的。
正在使用 v1.1 的人必须跟已经升级到 v1.2 的人一样，升级到 v1.2.1（或 v1.3）才能修复这个 bug。
当然应用开发者可以要求大版本号割裂为不同的应用，
例如 iOS 的 Reeder，v3 和 v4 就是两个独立的应用，买了 v3 的可以一直享受 v3 的升级，
即使是在 v4 发布一段时间后 v3 仍然享受独立的升级。

但大家可以看看，现在有多少开发者可以坚持服务于多个并存的版本号的？几乎没有…Reeder 一开始发布 v4 时保留了 v3，后来发现了一个严重的问题，最后这个问题就懒得在 v3 上修复了，直接把 v3 下架，让大家付费升级到 v4，否则无法修复这个严重的问题。
就算当年 Microsoft 可以保持多个 Windows 大版本并存，现在也不想搞了，直接就是 Windows 10 无限制升级下去。
理论上每半年 Windows 10 有一个大更新，但其实服务生命周期只有 18 个月（企业版 30 个月），也就是永远只对最近的 3 个大更新提供版本并存的升级服务。
好像 Windows XP 这种生命周期长大 12 年的产品早已不存在。

最后，按订阅收费也是现在的商业潮流。既然不再是按大版本购买收费了，那为什么还要允许多版本并存？所有订阅用户都在持续付费，都有权利及义务使用最新的版本。开发者不会再去维护并存的老版本。

尽管作者提议 GitHub Flow，但我觉得 GitHub Flow 仍然不完美。
有可能是我习惯了 Facebook 内部使用 Phabricator 的流程了，而现在我看到别人公司用 Phabricator 会觉得很顺手，反而自己的小项目用 GitHub 觉得有点烦。
（我懒得为自己的小项目搭建 Phabricator，所以只能用 GitHub Flow。）

GitHub Flow 最大的问题在于那个 merge。
首先，谁负责 merge？这在 GitHub 是个开放性的问题。
有时候是提交 Pull Request 的人负责，有时候是接受 Pull Request 的人负责。
这导致跟新的团队合作时责任很不明确。
尽管看起来点一下 merge 是谁都能做的事情，但其实背后有巨大的责任：
如果出现 merge conflict，谁负责手工解决？如果手工解决的结果导致持续集成不通过了，谁负责修复？
我喜欢 Facebook 的流程，因为保护 master 是第一优先级，任何时候 master 都要能够通过持续集成，能够持续发布到线上。
为此，规则被设计为没有人能够 merge，commit 必须是 rebase 后直接放到 master 的顶部。这时候责任就很明确了，冲突不可能发生在 merge 上，是谁写的 commit 就谁负责在自己那端 rebase 到 master 上。
既然是在自己那端做，如果存在 rebase conflict 那当然是写 commit 的人负责解决。我觉得这样的流程明显顺手很多！
此外 GitHub 还有一个问题，就是每次 Pull Request 都会导致新开一个 branch，但 Pull Request 通过后谁负责删除这个 branch 同样是不明确的。
现在 GitHub 可以设置为自动删除 Pull Request 用过的 branch 了，以前就会留下一大堆大家忘记删除的 branch，需要想办法删除。
（据说 AirBnb 专门有个脚本来负责删除 GitHub 上没用的远程 branch。）在这方面，Phabricator 就好很多。
它会让你配置一个专门用来做 code review 的 staging repo，然后每个 Diff（相当于 GitHub Pull Request）会以 tag 的形式发布到 staging repo，不会污染主 repo。
主 repo 看到的就是整洁的一系列 commit 在 master 上，没有 branch 也没有 merge commit。