package hard

import (
	"math"
)

/*
391. 完美矩形
给你一个数组 rectangles ，其中 rectangles[i] = [xi, yi, ai, bi] 表示一个坐标轴平行的矩形。这个矩形的左下顶点是 (xi, yi) ，右上顶点是 (ai, bi) 。

如果所有矩形一起精确覆盖了某个矩形区域，则返回 true ；否则，返回 false 。


示例 1：


输入：rectangles = [[1,1,3,3],[3,1,4,2],[3,2,4,4],[1,3,2,4],[2,3,3,4]]
输出：true
解释：5 个矩形一起可以精确地覆盖一个矩形区域。
示例 2：


输入：rectangles = [[1,1,2,3],[1,3,2,4],[3,1,4,2],[3,2,4,4]]
输出：false
解释：两个矩形之间有间隔，无法覆盖成一个矩形。
示例 3：


输入：rectangles = [[1,1,3,3],[3,1,4,2],[1,3,2,4],[3,2,4,4]]
输出：false
解释：图形顶端留有空缺，无法覆盖成一个矩形。
示例 4：


输入：rectangles = [[1,1,3,3],[3,1,4,2],[1,3,2,4],[2,2,4,4]]
输出：false
解释：因为中间有相交区域，虽然形成了矩形，但不是精确覆盖。


提示：

1 <= rectangles.length <= 2 * 104
rectangles[i].length == 4
-105 <= xi, yi, ai, bi <= 105

https://leetcode-cn.com/problems/perfect-rectangle/
*/

func isRectangleCover(rectangles [][]int) bool {
	points := make(map[point]struct{})
	var sumArea int
	for i := range rectangles {
		p := point{rectangles[i][0], rectangles[i][1]}
		if _, ok := points[p]; ok {
			delete(points, p)
		} else {
			points[p] = struct{}{}
		}
		p = point{rectangles[i][0], rectangles[i][3]}
		if _, ok := points[p]; ok {
			delete(points, p)
		} else {
			points[p] = struct{}{}
		}
		p = point{rectangles[i][2], rectangles[i][1]}
		if _, ok := points[p]; ok {
			delete(points, p)
		} else {
			points[p] = struct{}{}
		}
		p = point{rectangles[i][2], rectangles[i][3]}
		if _, ok := points[p]; ok {
			delete(points, p)
		} else {
			points[p] = struct{}{}
		}
		sumArea += (rectangles[i][2] - rectangles[i][0]) * (rectangles[i][3] - rectangles[i][1])
	}
	x2, y2 := math.MaxInt, math.MaxInt
	var x1, y1 int
	for p := range points {
		if p.x > x1 {
			x1 = p.x
		}
		if p.y > y1 {
			y1 = p.y
		}
		if p.x < x2 {
			x2 = p.x
		}
		if p.y < y2 {
			y2 = p.y
		}
	}
	return len(points) == 4 && (x1-x2)*(y1-y2) == sumArea
}
