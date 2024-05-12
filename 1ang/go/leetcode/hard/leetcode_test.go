package hard

import (
	"fmt"
	"test/leetcode"
	"testing"
)

func TestReverseKGroup(t *testing.T) {
	list := leetcode.NewList([]int{1, 2})
	fmt.Println(reverseKGroup(list, 2))
}

func TestFindMaximizedCapital(t *testing.T) {
	fmt.Println(findMaximizedCapital(10, 0, []int{1, 2, 3}, []int{0, 1, 2}))
}

func TestFullJustify(t *testing.T) {
	fmt.Println(fullJustify([]string{"Science", "is", "what", "we", "understand", "well", "enough", "to", "explain", "to", "a", "computer.", "Art", "is", "everything", "else", "we", "do"}, 20))
}

func TestFirstMissingPositive(t *testing.T) {
	fmt.Println(firstMissingPositive([]int{3, 4, -1, 1}))
}

func TestNumDecodings(t *testing.T) {
	fmt.Println(numDecodings("7*9*3*6*3*0*5*4*9*7*3*7*1*8*3*2*0*0*6*"))
}

func TestFindMinMoves(t *testing.T) {
	fmt.Println(findMinMoves2([]int{0, 3, 0}))
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

func TestRemoveInvalidParentheses(t *testing.T) {
	fmt.Println(removeInvalidParentheses("())(((()m)("))
}

func TestIsSelfCrossing(t *testing.T) {
	fmt.Println(isSelfCrossing2([]int{2, 1, 1, 2}))
}

func TestTrapRainWater(t *testing.T) {
	fmt.Println(trapRainWater([][]int{{12, 13, 1, 12}, {13, 4, 13, 12}, {13, 8, 10, 12}, {12, 13, 12, 12}, {13, 13, 13, 13}}))
}

func TestIsRectangleCover(t *testing.T) {
	fmt.Println(isRectangleCover([][]int{{0, 0, 4, 1}, {7, 0, 8, 2}, {5, 1, 6, 3}, {6, 0, 7, 2}, {4, 0, 5, 1}, {4, 2, 5, 3}, {2, 1, 4, 3}, {0, 2, 2, 3}, {0, 1, 2, 2}, {6, 2, 8, 3}, {5, 0, 6, 1}, {4, 1, 5, 2}}))
}

func TestMaxSumOfThreeSubarrays(t *testing.T) {
	fmt.Println(maxSumOfThreeSubarrays3([]int{17, 7, 19, 11, 1, 19, 17, 6, 13, 18, 2, 7, 12, 16, 16, 18, 9, 3, 19, 5}, 6))
}
