package middle

/*
36. 有效的数独
请你判断一个 9x9 的数独是否有效。只需要 根据以下规则 ，验证已经填入的数字是否有效即可。

数字 1-9 在每一行只能出现一次。
数字 1-9 在每一列只能出现一次。
数字 1-9 在每一个以粗实线分隔的 3x3 宫内只能出现一次。（请参考示例图）
数独部分空格内已填入了数字，空白格用 '.' 表示。

注意：

一个有效的数独（部分已被填充）不一定是可解的。
只需要根据以上规则，验证已经填入的数字是否有效即可。


示例 1：


输入：board =
[["5","3",".",".","7",".",".",".","."]
,["6",".",".","1","9","5",".",".","."]
,[".","9","8",".",".",".",".","6","."]
,["8",".",".",".","6",".",".",".","3"]
,["4",".",".","8",".","3",".",".","1"]
,["7",".",".",".","2",".",".",".","6"]
,[".","6",".",".",".",".","2","8","."]
,[".",".",".","4","1","9",".",".","5"]
,[".",".",".",".","8",".",".","7","9"]]
输出：true
示例 2：

输入：board =
[["8","3",".",".","7",".",".",".","."]
,["6",".",".","1","9","5",".",".","."]
,[".","9","8",".",".",".",".","6","."]
,["8",".",".",".","6",".",".",".","3"]
,["4",".",".","8",".","3",".",".","1"]
,["7",".",".",".","2",".",".",".","6"]
,[".","6",".",".",".",".","2","8","."]
,[".",".",".","4","1","9",".",".","5"]
,[".",".",".",".","8",".",".","7","9"]]
输出：false
解释：除了第一行的第一个数字从 5 改为 8 以外，空格内其他数字均与 示例1 相同。 但由于位于左上角的 3x3 宫内有两个 8 存在, 因此这个数独是无效的。


提示：

board.length == 9
board[i].length == 9
board[i][j] 是一位数字或者 '.'

https://leetcode-cn.com/problems/valid-sudoku/
*/

func isValidSudoku(board [][]byte) bool {
	efficient := make([][]bool, 27)
	for i, b := range board {
		efficient[i] = make([]bool, 9)
		for j, c := range b {
			if c == '.' {
				continue
			}
			if efficient[i][c-'0'] {
				return false
			} else {
				efficient[i][c-'0'] = true
			}
			idx := j + 9
			if efficient[idx] == nil {
				efficient[idx] = make([]bool, 27)
			}
			if efficient[idx][c-'0'] {
				return false
			} else {
				efficient[idx][c-'0'] = true
			}
			idx = (i / 3 * 3) + j/3 + 18
			if efficient[idx] == nil {
				efficient[idx] = make([]bool, 27)
			}
			if efficient[idx][c-'0'] {
				return false
			} else {
				efficient[idx][c-'0'] = true
			}
		}
	}
	return true
}

func isValidSudoku2(board [][]byte) bool {
	validArr1 := make([]byte, 9)
	validArr2 := make([]byte, 9)
	validArr3 := make([]byte, 9)
	for i, b := range board {
		num := byte(i + 1)
		for j := range b {
			if board[i][j] != '.' {
				idx := board[i][j] - '1'
				if validArr1[idx] != num {
					validArr1[idx] = num
				} else {
					return false
				}
			}
			if board[j][i] != '.' {
				idx := board[j][i] - '1'
				if validArr2[idx] != num {
					validArr2[idx] = num
				} else {
					return false
				}
			}
			x := j/3 + (i/3)*3
			y := j%3 + (i%3)*3
			if board[x][y] != '.' {
				idx := board[x][y] - '1'
				if validArr3[idx] != num {
					validArr3[idx] = num
				} else {
					return false
				}
			}
		}
	}
	return true
}
