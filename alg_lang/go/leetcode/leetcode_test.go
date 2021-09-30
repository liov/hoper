package leetcode

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	println(isValid("()()()"))
}

func TestCompareVersion(t *testing.T) {
	println(compareVersion2("11", "10"))
}

func TestLongestPalindrome(t *testing.T) {
	println(longestPalindrome("ac"))
}

func TestConvert(t *testing.T) {
	println(convert("PAYPALISHIRING", 3))
}

func TestReverseKGroup(t *testing.T) {
	list := NewList([]int{1, 2})
	fmt.Println(reverseKGroup(list, 2))
}

func TestMyAtoi(t *testing.T) {
	fmt.Println(myAtoi("00000-42a1234"))
}

func TestSearch2(t *testing.T) {
	fmt.Println(search2([]int{8, 1, 2, 3, 4, 5, 6, 7}, 6))
}

func TestFindMaximizedCapital(t *testing.T) {
	fmt.Println(findMaximizedCapital(10, 0, []int{1, 2, 3}, []int{0, 1, 2}))
}

func TestFullJustify(t *testing.T) {
	fmt.Println(fullJustify([]string{"Science", "is", "what", "we", "understand", "well", "enough", "to", "explain", "to", "a", "computer.", "Art", "is", "everything", "else", "we", "do"}, 20))
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

func TestFirstMissingPositive(t *testing.T) {
	fmt.Println(firstMissingPositive([]int{3, 4, -1, 1}))
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

func TestNumDecodings(t *testing.T) {
	fmt.Println(numDecodings("7*9*3*6*3*0*5*4*9*7*3*7*1*8*3*2*0*0*6*"))
}

func TestPathSum(t *testing.T) {
	tree := NewTree([]int{10, 5, -3, 3, 2, 0, 11, 3, -2, 0, 1})
	fmt.Println(pathSum(tree, 8))
}

func TestFindMinMoves(t *testing.T) {
	fmt.Println(findMinMoves2([]int{0, 3, 0}))
}

func TestComputeArea(t *testing.T) {
	fmt.Println(computeArea(-2, -2, 2, 2, -2, -2, 2, 2))
}
