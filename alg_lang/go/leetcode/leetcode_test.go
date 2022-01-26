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

func TestFindRepeatedDnaSequences(t *testing.T) {
	fmt.Println(findRepeatedDnaSequences("AAAAAAAAAAA"))
}

func TestSummaryRanges(t *testing.T) {
	obj := Constructor()
	set := []int{49, 97, 53, 5, 33, 65, 62, 51, 100, 38, 61, 45, 74, 27, 64, 17, 36, 17, 96, 12, 79, 32, 68, 90, 77, 18, 39, 12, 93, 9, 87, 42, 60, 71, 12, 45, 55, 40, 78, 81, 26, 70, 61, 56, 66, 33, 7, 70, 1, 11, 92, 51, 90, 100, 85, 80, 0, 78, 63, 42, 31, 93, 41, 90, 8, 24, 72, 28, 30, 18, 69, 57, 11, 10, 40, 65, 62, 13, 38, 70, 37, 90, 15, 70, 42, 69, 26, 77, 70, 75, 36, 56, 11, 76, 49, 40, 73, 30, 37, 23}
	for i := range set {
		obj.AddNum(set[i])
		fmt.Println(obj.GetIntervals())
	}
}

func TestNumberToWords(t *testing.T) {
	fmt.Println(numberToWords(0))
}

func TestDivide(t *testing.T) {
	fmt.Println(divide(7, -3))
}

func TestSearchMatrix(t *testing.T) {
	fmt.Println(searchMatrix([][]int{}, 1))
}

func TestRemoveInvalidParentheses(t *testing.T) {
	fmt.Println(removeInvalidParentheses("())(((()m)("))
}

func TestIsSelfCrossing(t *testing.T) {
	fmt.Println(isSelfCrossing2([]int{2, 1, 1, 2}))
}

func TestSingleNumber(t *testing.T) {
	fmt.Println(singleNumber([]int{2, 1, 1, 2, 5, 6}))
}

func TestTrapRainWater(t *testing.T) {
	fmt.Println(trapRainWater([][]int{{12, 13, 1, 12}, {13, 4, 13, 12}, {13, 8, 10, 12}, {12, 13, 12, 12}, {13, 13, 13, 13}}))
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

func TestIsRectangleCover(t *testing.T) {
	fmt.Println(isRectangleCover([][]int{{0, 0, 4, 1}, {7, 0, 8, 2}, {5, 1, 6, 3}, {6, 0, 7, 2}, {4, 0, 5, 1}, {4, 2, 5, 3}, {2, 1, 4, 3}, {0, 2, 2, 3}, {0, 1, 2, 2}, {6, 2, 8, 3}, {5, 0, 6, 1}, {4, 1, 5, 2}}))
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

func TestLargestSumAfterKNegations(t *testing.T) {
	fmt.Println(largestSumAfterKNegations([]int{2, -3, -1, 5, -4}, 2))
}

func TestMaxSumOfThreeSubarrays(t *testing.T) {
	fmt.Println(maxSumOfThreeSubarrays3([]int{17, 7, 19, 11, 1, 19, 17, 6, 13, 18, 2, 7, 12, 16, 16, 18, 9, 3, 19, 5}, 6))
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
