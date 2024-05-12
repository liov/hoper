package middle

import (
	"fmt"
	"testing"
)

func TestCompareVersion(t *testing.T) {
	println(compareVersion2("11", "10"))
}

func TestLongestPalindrome(t *testing.T) {
	println(longestPalindrome("ac"))
}

func TestConvert(t *testing.T) {
	println(convert("PAYPALISHIRING", 3))
}

func TestMyAtoi(t *testing.T) {
	fmt.Println(myAtoi("00000-42a1234"))
}

func TestSearch2(t *testing.T) {
	fmt.Println(search2([]int{8, 1, 2, 3, 4, 5, 6, 7}, 6))
}

func TestSwapPairs(t *testing.T) {
	list := NewList([]int{1, 2, 3})
	fmt.Println(swapPairs(list))
}

func TestNextPermutation(t *testing.T) {
	nums := []int{1, 2, 3}
	nextPermutation(nums)
	fmt.Println(nums)
}

func TestNumberOfBoomerangs(t *testing.T) {
	fmt.Println(numberOfBoomerangs([][]int{{0, 0}, {1, 0}, {2, 0}}))
}

func TestIsBipartite(t *testing.T) {
	fmt.Println(isBipartite([][]int{{1, 3}, {0, 2}, {1, 3}, {0, 2}}))
}

func TestMaxArea(t *testing.T) {
	fmt.Println(maxArea([]int{2, 3, 4, 5, 18, 17, 6}))
}

func TestFlatten(t *testing.T) {
	nodes := make([]*Node, 12)
	for i := 0; i < len(nodes); i++ {
		nodes[i] = &Node{Val: i + 1}
	}
	link(nodes, 0, 1)
	link(nodes, 1, 2)
	link(nodes, 2, 3)
	link(nodes, 3, 4)
	link(nodes, 4, 5)
	nodes[2].Child = nodes[6]
	link(nodes, 6, 7)
	link(nodes, 7, 8)
	link(nodes, 8, 9)
	nodes[7].Child = nodes[10]
	link(nodes, 10, 11)
	fmt.Println(flatten(nodes[0]))
}

func link(nodes []*Node, i, j int) {
	nodes[i].Next = nodes[j]
	nodes[j].Prev = nodes[i]
}

func TestGetSum(t *testing.T) {
	fmt.Println(getSum(-2, -3))
}

func TestPathSum(t *testing.T) {
	tree := NewTree([]int{10, 5, -3, 3, 2, 0, 11, 3, -2, 0, 1})
	fmt.Println(pathSum((tree), 8))
}

func TestComputeArea(t *testing.T) {
	fmt.Println(computeArea(-2, -2, 2, 2, -2, -2, 2, 2))
}

func TestFindRepeatedDnaSequences(t *testing.T) {
	fmt.Println(findRepeatedDnaSequences("AAAAAAAAAAA"))
}

func TestDivide(t *testing.T) {
	fmt.Println(divide(7, -3))
}

func TestSearchMatrix(t *testing.T) {
	fmt.Println(searchMatrix([][]int{}, 1))
}

func TestSingleNumber(t *testing.T) {
	fmt.Println(singleNumber([]int{2, 1, 1, 2, 5, 6}))
}

func TestLongestSubsequence(t *testing.T) {
	fmt.Println(longestSubsequence([]int{1, 2, 3, 4}, 1))
}

func TestGetMoneyAmount(t *testing.T) {
	fmt.Println(getMoneyAmount2(16))
}

func TestBulbSwitch(t *testing.T) {
	fmt.Println(bulbSwitch(16))
}

func TestMaxProduct(t *testing.T) {
	fmt.Println(maxProduct([]string{"a", "aa", "aaa", "aaaa"}))
}

func TestOriginalDigits(t *testing.T) {
	fmt.Println(originalDigits("zeroonetwothreefourfivesixseveneightnine"))
}

func TestFindNthDigit(t *testing.T) {
	fmt.Println(findNthDigit(500))
}

func TestValidTicTacToe(t *testing.T) {
	fmt.Println(validTicTacToe([]string{"XOX", " X ", "   "}))
}

func TestRepeatedStringMatch(t *testing.T) {
	fmt.Println(repeatedStringMatch("abcd", "cdabcdacdabcda"))
}

func TestIsAdditiveNumber(t *testing.T) {
	fmt.Println(isAdditiveNumber("199111992"))
}

func TestFindMinDifference(t *testing.T) {
	fmt.Println(findMinDifference([]string{"01:01", "02:01", "03:00"}))
}

func TestKnightProbability(t *testing.T) {
	fmt.Println(knightProbability(10, 13, 5, 3))
}

func TestSearchRange(t *testing.T) {
	t.Log(searchRange([]int{5, 7, 7, 8, 8, 10}, 8))
}

func TestRob(t *testing.T) {
	t.Log(rob([]int{2, 7, 9, 3, 1}))
}

func TestDecodeString(t *testing.T) {
	t.Log(decodeString("3[a]2[bc]"))
}

func TestCoinChange(t *testing.T) {
	t.Log(coinChange([]int{1, 2, 5}, 11))
}

func TestWordBreak(t *testing.T) {
	t.Log(wordBreak("leetcode", []string{"leet", "code"}))
}

func TestPartition(t *testing.T) {
	t.Log(partition("aabaa"))
}

func TestPermute(t *testing.T) {
	t.Log(permute([]int{1, 2, 3, 4}))
}

func TestSubsets(t *testing.T) {
	t.Log(subsets([]int{1, 2, 3}))
}

func TestGenerateParenthesis(t *testing.T) {
	t.Log(generateParenthesis(3))
}
