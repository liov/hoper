package middle

/*
240. 搜索二维矩阵 II
编写一个高效的算法来搜索 m x n 矩阵 matrix 中的一个目标值 target 。该矩阵具有以下特性：

每行的元素从左到右升序排列。
每列的元素从上到下升序排列。


示例 1：


输入：matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 5
输出：true
示例 2：


输入：matrix = [[1,4,7,11,15],[2,5,8,12,19],[3,6,9,16,22],[10,13,14,17,24],[18,21,23,26,30]], target = 20
输出：false


提示：

m == matrix.length
n == matrix[i].length
1 <= n, m <= 300
-10^9 <= matrix[i][j] <= 10^9
每行的所有元素从左到右升序排列
每列的所有元素从上到下升序排列
-10^9 <= target <= 10^9

https://leetcode-cn.com/problems/search-a-2d-matrix-ii/
*/

func searchMatrix(matrix [][]int, target int) bool {
	rl, rr := 0, len(matrix[0])-1
	cl, cr := 0, len(matrix)-1

	if searchMatrixHelper(matrix, target, rl, rr, cl, cr) {
		return true
	}
	return false
}

func searchMatrixHelper(matrix [][]int, target int, rl, rr, cl, cr int) bool {
	x1, x2, y1, y2 := rl, rr, cl, cr
	if target > matrix[cr][rr] || target < matrix[cl][rl] {
		return false
	}

	for {
		rm := (rl + rr) >> 1
		cm := (cl + cr) >> 1
		if matrix[cm][rm] == target {
			return true
		}
		if matrix[cm][rm] > target {
			if rm > x1 {
				rr = rm - 1
			}
			if cm > y1 {
				cr = cm - 1
			}
			if matrix[cr][rr] < target {
				l, r := y1, cr
				if rm <= x2 && matrix[l][rm] <= target {
					for l <= r {
						m := (l + r) >> 1
						if matrix[m][rm] == target {
							return true
						}
						if matrix[m][rm] < target {
							l = m + 1
							continue
						}
						r = m - 1
					}
				}
				l, r = x1, rr
				if cm <= y2 && matrix[cm][l] <= target {
					for l <= r {
						m := (l + r) >> 1
						if matrix[cm][m] == target {
							return true
						}
						if matrix[cm][m] < target {
							l = m + 1
							continue
						}
						r = m - 1
					}
				}
				if rm < x2 && cm > y1 && searchMatrixHelper(matrix, target, rm+1, x2, y1, cm-1) {
					return true
				}

				if cm < y2 && rm > x1 && searchMatrixHelper(matrix, target, x1, rm-1, cm+1, y2) {
					return true
				}
				return false
			}
			continue
		}
		if rm == rr && cm == cr {
			return false
		}
		if rm < x2 {
			rl = rm + 1
		}
		if cm < y2 {
			cl = cm + 1
		}
	}
}

func searchMatrix2(matrix [][]int, target int) bool {
	m, n := len(matrix), len(matrix[0])
	x, y := 0, n-1
	for x < m && y >= 0 {
		if matrix[x][y] == target {
			return true
		}
		if matrix[x][y] > target {
			y--
		} else {
			x++
		}
	}
	return false
}
