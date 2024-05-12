package middle

/*
794. 有效的井字游戏
给你一个字符串数组 board 表示井字游戏的棋盘。当且仅当在井字游戏过程中，棋盘有可能达到 board 所显示的状态时，才返回 true 。

井字游戏的棋盘是一个 3 x 3 数组，由字符 ' '，'X' 和 'O' 组成。字符 ' ' 代表一个空位。

以下是井字游戏的规则：

玩家轮流将字符放入空位（' '）中。
玩家 1 总是放字符 'X' ，而玩家 2 总是放字符 'O' 。
'X' 和 'O' 只允许放置在空位中，不允许对已放有字符的位置进行填充。
当有 3 个相同（且非空）的字符填充任何行、列或对角线时，游戏结束。
当所有位置非空时，也算为游戏结束。
如果游戏结束，玩家不允许再放置字符。


示例 1：


输入：board = ["O  ","   ","   "]
输出：false
解释：玩家 1 总是放字符 "X" 。
示例 2：


输入：board = ["XOX"," X ","   "]
输出：false
解释：玩家应该轮流放字符。
示例 3：


输入：board = ["XXX","   ","OOO"]
输出：false
Example 4:


输入：board = ["XOX","O O","XOX"]
输出：true


提示：

board.length == 3
board[i].length == 3
board[i][j] 为 'X'、'O' 或 ' '

https://leetcode-cn.com/problems/valid-tic-tac-toe-state/
*/
//1、X必须比O多1个或者相同，否则false； 2、X和O不能同时都是赢家，否则false； 3、X赢的时候必须比O多1个，否则false； 4、O赢得时候个数与X相同，否则false；
func validTicTacToe(board []string) bool {
	var xcount, ocount int
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 'X' {
				xcount++
			}
			if board[i][j] == 'O' {
				ocount++
			}
		}
	}
	if win(board, 'X') && win(board, 'O') {
		return false
	}
	if win(board, 'X') && xcount-ocount != 1 {
		return false
	}
	if win(board, 'O') && xcount != ocount {
		return false
	}
	if ocount > xcount || xcount-ocount > 1 {
		return false
	}
	return true
}

func win(board []string, b byte) bool {
	if board[0][0] == b && board[1][0] == b && board[2][0] == b {
		return true
	}
	if board[0][0] == b && board[0][1] == b && board[0][2] == b {
		return true
	}
	if board[0][0] == b && board[1][1] == b && board[2][2] == b {
		return true
	}
	if board[0][1] == b && board[1][1] == b && board[2][1] == b {
		return true
	}
	if board[1][0] == b && board[1][1] == b && board[1][2] == b {
		return true
	}
	if board[0][2] == b && board[1][2] == b && board[2][2] == b {
		return true
	}
	if board[0][2] == b && board[1][2] == b && board[2][2] == b {
		return true
	}
	if board[0][2] == b && board[1][1] == b && board[2][0] == b {
		return true
	}
	if board[2][0] == b && board[2][1] == b && board[2][2] == b {
		return true
	}
	return false
}
