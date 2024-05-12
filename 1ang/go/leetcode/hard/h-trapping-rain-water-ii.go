package hard

import (
	"container/heap"
	"test/leetcode"
)

/*
407. 接雨水 II

给你一个 m x n 的矩阵，其中的值均为非负整数，代表二维高度图每个单元的高度，请计算图中形状最多能接多少体积的雨水。



示例 1:



输入: heightMap = [[1,4,3,1,3,2],[3,2,1,3,2,4],[2,3,3,2,3,1]]
输出: 4
解释: 下雨后，雨水将会被上图蓝色的方块中。总的接雨水量为1+2+1=4。
示例 2:



输入: heightMap = [[3,3,3,3,3],[3,2,2,2,3],[3,2,1,2,3],[3,2,2,2,3],[3,3,3,3,3]]
输出: 10


提示:

m == heightMap.length
n == heightMap[i].length
1 <= m, n <= 200
0 <= heightMap[i][j] <= 2 * 10^4

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/trapping-rain-water-ii
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/

// 围边战术,超时
func trapRainWater(heightMap [][]int) int {
	m, n := len(heightMap), len(heightMap[0])
	if m < 3 || n < 3 {
		return 0
	}
	drop := make([][]bool, m)
	for i := range drop {
		drop[i] = make([]bool, n)
	}
	side := make(hp, 0)
	for i, row := range heightMap {
		for j, h := range row {
			if i == 0 || i == m-1 || j == 0 || j == n-1 {
				heap.Push(&side, cell{h, i, j})
				drop[i][j] = true
			}
		}
	}
	round := [][]int{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
	var x, y, result int
	for side.Len() > 0 {
		p1 := heap.Pop(&side).(cell)
		for _, p := range round {
			x = p1.x + p[0]
			y = p1.y + p[1]

			if x <= 0 || x >= m-1 || y <= 0 || y >= n-1 {
				continue
			}
			if drop[x][y] {
				continue
			}
			drop[x][y] = true
			if heightMap[x][y] <= p1.h {
				result += p1.h - heightMap[x][y]
				heap.Push(&side, cell{p1.h, x, y})
			} else {
				heap.Push(&side, cell{heightMap[x][y], x, y})
			}
		}

	}
	return result
}

func trapRainWater2(heightMap [][]int) (ans int) {
	m, n := len(heightMap), len(heightMap[0])
	maxHeight := 0
	for _, row := range heightMap {
		for _, h := range row {
			maxHeight = leetcode.max(maxHeight, h)
		}
	}

	water := make([][]int, m)
	for i := range water {
		water[i] = make([]int, n)
		for j := range water[i] {
			water[i][j] = maxHeight
		}
	}
	type pair struct{ x, y int }
	q := []pair{}
	for i, row := range heightMap {
		for j, h := range row {
			if (i == 0 || i == m-1 || j == 0 || j == n-1) && h < water[i][j] {
				water[i][j] = h
				q = append(q, pair{i, j})
			}
		}
	}

	dirs := []int{-1, 0, 1, 0, -1}
	for len(q) > 0 {
		p := q[0]
		q = q[1:]
		x, y := p.x, p.y
		for i := 0; i < 4; i++ {
			nx, ny := x+dirs[i], y+dirs[i+1]
			if 0 <= nx && nx < m && 0 <= ny && ny < n && water[nx][ny] > water[x][y] && water[nx][ny] > heightMap[nx][ny] {
				water[nx][ny] = leetcode.max(water[x][y], heightMap[nx][ny])
				q = append(q, pair{nx, ny})
			}
		}
	}

	for i, row := range heightMap {
		for j, h := range row {
			ans += water[i][j] - h
		}
	}
	return
}

func trapRainWater3(heightMap [][]int) (ans int) {
	m, n := len(heightMap), len(heightMap[0])
	if m <= 2 || n <= 2 {
		return
	}

	vis := make([][]bool, m)
	for i := range vis {
		vis[i] = make([]bool, n)
	}
	h := &hp{}
	for i, row := range heightMap {
		for j, v := range row {
			if i == 0 || i == m-1 || j == 0 || j == n-1 {
				heap.Push(h, cell{v, i, j})
				vis[i][j] = true
			}
		}
	}

	dirs := []int{-1, 0, 1, 0, -1}
	for h.Len() > 0 {
		cur := heap.Pop(h).(cell)
		for k := 0; k < 4; k++ {
			nx, ny := cur.x+dirs[k], cur.y+dirs[k+1]
			if 0 <= nx && nx < m && 0 <= ny && ny < n && !vis[nx][ny] {
				if heightMap[nx][ny] < cur.h {
					ans += cur.h - heightMap[nx][ny]
				}
				vis[nx][ny] = true
				heap.Push(h, cell{leetcode.max(heightMap[nx][ny], cur.h), nx, ny})
			}
		}
	}
	return
}

type cell struct{ h, x, y int }
type hp []cell

func (h hp) Len() int            { return len(h) }
func (h hp) Less(i, j int) bool  { return h[i].h < h[j].h }
func (h hp) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(v interface{}) { *h = append(*h, v.(cell)) }
func (h *hp) Pop() interface{}   { a := *h; v := a[len(a)-1]; *h = a[:len(a)-1]; return v }
