package middle

/*
447. 回旋镖的数量
给定平面上 n 对 互不相同 的点 points ，其中 points[i] = [xi, yi] 。回旋镖 是由点 (i, j, k) 表示的元组 ，其中 i 和 j 之间的距离和 i 和 k 之间的距离相等（需要考虑元组的顺序）。

返回平面上所有回旋镖的数量。

示例 1：

输入：points = [[0,0],[1,0],[2,0]]
输出：2
解释：两个回旋镖为 [[1,0],[0,0],[2,0]] 和 [[1,0],[2,0],[0,0]]
示例 2：

输入：points = [[1,1],[2,2],[3,3]]
输出：2
示例 3：

输入：points = [[1,1]]
输出：0

提示：

n == points.length
1 <= n <= 500
points[i].length == 2
-104 <= xi, yi <= 104
所有点都 互不相同

https://leetcode-cn.com/problems/number-of-boomerangs
*/
func numberOfBoomerangs(points [][]int) int {
	if len(points) < 3 {
		return 0
	}
	var ret int
	for j := 0; j < len(points); j++ {
		points[j], points[0] = points[0], points[j]
		var m = make(map[int]int)
		for i := 1; i < len(points); i++ {
			limit := (points[i][0]-points[0][0])*(points[i][0]-points[0][0]) + (points[i][1]-points[0][1])*(points[i][1]-points[0][1])
			m[limit] = m[limit] + 1
		}
		for _, count := range m {
			ret += count * (count - 1)
		}
		points[j], points[0] = points[0], points[j]
	}
	return ret
}
