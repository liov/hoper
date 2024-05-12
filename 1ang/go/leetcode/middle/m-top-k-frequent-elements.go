package middle

/*
347. 前 K 个高频元素
中等
给你一个整数数组 nums 和一个整数 k ，请你返回其中出现频率前 k 高的元素。你可以按 任意顺序 返回答案。



示例 1:

输入: nums = [1,1,1,2,2,3], k = 2
输出: [1,2]
示例 2:

输入: nums = [1], k = 1
输出: [1]


提示：

1 <= nums.length <= 105
k 的取值范围是 [1, 数组中不相同的元素的个数]
题目数据保证答案唯一，换句话说，数组中前 k 个高频元素的集合是唯一的


进阶：你所设计算法的时间复杂度 必须 优于 O(n log n) ，其中 n 是数组大小。
*/

func topKFrequent(nums []int, k int) []int {
	numMap := make(map[int]int)
	for _, num := range nums {
		numMap[num]++
	}
	heap := make([]NumTimes, k+1)
	for num, times := range numMap {
		heap[0].Num = num
		heap[0].Times = times
		down(heap, 0)
	}
	var ans []int
	for i := 1; i < len(heap); i++ {
		ans = append(ans, heap[i].Num)
	}
	return ans
}

type NumTimes struct {
	Num, Times int
}

func down(nums []NumTimes, i int) {
	child := i*2 + 1
	for child < len(nums) {
		if child+1 < len(nums) && nums[child].Times > nums[child+1].Times {
			child++
		}
		if nums[child].Times > nums[i].Times {
			break
		}
		nums[i], nums[child] = nums[child], nums[i]
		i = child
		child = child*2 + 1
	}
}
