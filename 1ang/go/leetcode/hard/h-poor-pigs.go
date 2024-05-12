package hard

import "math"

/*
458. 可怜的小猪
有 buckets 桶液体，其中 正好 有一桶含有毒药，其余装的都是水。它们从外观看起来都一样。为了弄清楚哪只水桶含有毒药，你可以喂一些猪喝，通过观察猪是否会死进行判断。不幸的是，你只有 minutesToTest 分钟时间来确定哪桶液体是有毒的。

喂猪的规则如下：

选择若干活猪进行喂养
可以允许小猪同时饮用任意数量的桶中的水，并且该过程不需要时间。
小猪喝完水后，必须有 minutesToDie 分钟的冷却时间。在这段时间里，你只能观察，而不允许继续喂猪。
过了 minutesToDie 分钟后，所有喝到毒药的猪都会死去，其他所有猪都会活下来。
重复这一过程，直到时间用完。
给你桶的数目 buckets ，minutesToDie 和 minutesToTest ，返回在规定时间内判断哪个桶有毒所需的 最小 猪数。



示例 1：

输入：buckets = 1000, minutesToDie = 15, minutesToTest = 60
输出：5
示例 2：

输入：buckets = 4, minutesToDie = 15, minutesToTest = 15
输出：2
示例 3：

输入：buckets = 4, minutesToDie = 15, minutesToTest = 30
输出：2


提示：

1 <= buckets <= 1000
1 <= minutesToDie <= minutesToTest <= 100

https://leetcode-cn.com/problems/poor-pigs/
*/

/*
其实香农已经在《信息论》（信息熵）中给过我们结论了——我们一共可以进行n轮实验（n = minutesToTest / minutesToDie）：

经过所有实验，一只小猪能有多少种状态？第一轮就死、第二轮死、...、第n轮死，以及生还，所以一共有n + 1种状态
n + 1种状态所携带的信息为log_2(n + 1)比特，这也是一只小猪最多提供的信息量
而”buckets瓶液体中哪一瓶是毒“这件事，也有buckets种可能性，所以需要的信息量是log_2(buckets)
注：以上所有事件、状态都是先验等概的，所以可以直接对2取对数得到信息熵

因此一定存在一种“合理设计”的实验，使得我们只要有k只猪猪：满足 k * log_2(n + 1) >= log_2(buckets)时，则我们一定能得到足够的信息量去判断哪一瓶是毒。
*/

/*
这里以一个小时，1000桶举例举例，一头猪可以获得的信息有 0-15分钟死/15-30分钟死/30-45分钟死/45-60分钟死/不死， 一共五个状态，可以对应五进制的0/1/2/3/4，这头猪在实验完成后的状态就对应一个5进制数。选择4头猪可以提供5^4=625种，选择5头猪可以提供5^5=3125种，所以这种情况要5头猪。

总的来说，总时间/冷却时间+1=轮次，这里的轮次就是一头猪可以提供的信息，m头猪在n轮次情况下可以提供n^m个信息。
*/
func poorPigs(buckets int, minutesToDie int, minutesToTest int) int {
	states := minutesToTest/minutesToDie + 1
	return int(math.Ceil(math.Log(float64(buckets)) / math.Log(float64(states))))
}
