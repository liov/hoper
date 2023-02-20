package leetcode

import "sort"

/*
539. 最小时间差
给定一个 24 小时制（小时:分钟 "HH:MM"）的时间列表，找出列表中任意两个时间的最小时间差并以分钟数表示。



示例 1：

输入：timePoints = ["23:59","00:00"]
输出：1
示例 2：

输入：timePoints = ["00:00","23:59","00:00"]
输出：0


提示：

2 <= timePoints <= 2 * 10^4
timePoints[i] 格式为 "HH:MM"

https://leetcode-cn.com/problems/minimum-time-difference/
*/

func findMinDifference(timePoints []string) int {
	if len(timePoints) >= 1440 {
		return 0
	}
	var time []int
	for _, timePoint := range timePoints {
		m := (int(timePoint[0]-'0')*10+int(timePoint[1]-'0'))*60 + int(timePoint[3]-'0')*10 + int(timePoint[4]-'0')
		time = append(time, m)
	}
	sort.Ints(time)
	min := 1440
	for i := 1; i < len(time); i++ {
		if time[i]-time[i-1] < min {
			min = time[i] - time[i-1]
		}
		if min == 0 {
			return 0
		}
	}
	if min > 1440-time[len(time)-1]+time[0] {
		min = 1440 - time[len(time)-1] + time[0]
	}
	return min
}
