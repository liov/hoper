package middle

/*
223. 矩形面积
给你 二维 平面上两个 由直线构成的 矩形，请你计算并返回两个矩形覆盖的总面积。

每个矩形由其 左下 顶点和 右上 顶点坐标表示：

第一个矩形由其左下顶点 (ax1, ay1) 和右上顶点 (ax2, ay2) 定义。
第二个矩形由其左下顶点 (bx1, by1) 和右上顶点 (bx2, by2) 定义。


示例 1：

Rectangle Area
输入：ax1 = -3, ay1 = 0, ax2 = 3, ay2 = 4, bx1 = 0, by1 = -1, bx2 = 9, by2 = 2
输出：45
示例 2：

输入：ax1 = -2, ay1 = -2, ax2 = 2, ay2 = 2, bx1 = -2, by1 = -2, bx2 = 2, by2 = 2
输出：16


提示：

-10^4 <= ax1, ay1, ax2, ay2, bx1, by1, bx2, by2 <= 10^4

https://leetcode-cn.com/problems/rectangle-area/
*/

func computeArea(ax1 int, ay1 int, ax2 int, ay2 int, bx1 int, by1 int, bx2 int, by2 int) int {
	sum := (ax2-ax1)*(ay2-ay1) + (bx2-bx1)*(by2-by1)
	var cx1, cy1, cx2, cy2 int
	if bx1 >= ax1 && bx1 <= ax2 {
		cx1 = bx1
	} else if ax1 >= bx1 && ax1 <= bx2 {
		cx1 = ax1
	}
	if bx2 >= ax1 && bx2 <= ax2 {
		cx2 = bx2
	} else if ax2 >= bx1 && ax2 <= bx2 {
		cx2 = ax2
	}
	if by1 >= ay1 && by1 <= ay2 {
		cy1 = by1
	} else if ay1 >= by1 && ay1 <= by2 {
		cy1 = ay1
	}
	if by2 >= ay1 && by2 <= ay2 {
		cy2 = by2
	} else if ay2 >= by1 && ay2 <= by2 {
		cy2 = ay2
	}
	return sum - (cx2-cx1)*(cy2-cy1)
}
