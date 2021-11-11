package leetcode

/*
352. 将数据流变为多个不相交区间
 给你一个由非负整数 a1, a2, ..., an 组成的数据流输入，请你将到目前为止看到的数字总结为不相交的区间列表。

实现 SummaryRanges 类：

SummaryRanges() 使用一个空数据流初始化对象。
void addNum(int val) 向数据流中加入整数 val 。
int[][] getIntervals() 以不相交区间 [starti, endi] 的列表形式返回对数据流中整数的总结。


示例：

输入：
["SummaryRanges", "addNum", "getIntervals", "addNum", "getIntervals", "addNum", "getIntervals", "addNum", "getIntervals", "addNum", "getIntervals"]
[[], [1], [], [3], [], [7], [], [2], [], [6], []]
输出：
[null, null, [[1, 1]], null, [[1, 1], [3, 3]], null, [[1, 1], [3, 3], [7, 7]], null, [[1, 3], [7, 7]], null, [[1, 3], [6, 7]]]

解释：
SummaryRanges summaryRanges = new SummaryRanges();
summaryRanges.addNum(1);      // arr = [1]
summaryRanges.getIntervals(); // 返回 [[1, 1]]
summaryRanges.addNum(3);      // arr = [1, 3]
summaryRanges.getIntervals(); // 返回 [[1, 1], [3, 3]]
summaryRanges.addNum(7);      // arr = [1, 3, 7]
summaryRanges.getIntervals(); // 返回 [[1, 1], [3, 3], [7, 7]]
summaryRanges.addNum(2);      // arr = [1, 2, 3, 7]
summaryRanges.getIntervals(); // 返回 [[1, 3], [7, 7]]
summaryRanges.addNum(6);      // arr = [1, 2, 3, 6, 7]
summaryRanges.getIntervals(); // 返回 [[1, 3], [6, 7]]


提示：

0 <= val <= 10^4
最多调用 addNum 和 getIntervals 方法 3 * 10^4 次


进阶：如果存在大量合并，并且与数据流的大小相比，不相交区间的数量很小，该怎么办?

https://leetcode-cn.com/problems/data-stream-as-disjoint-intervals/
*/

type SummaryRanges struct {
	ret [][]int
	arr []bool
}

func Constructor() SummaryRanges {
	return SummaryRanges{arr: make([]bool, 10e4)}
}

func (this *SummaryRanges) AddNum(val int) {
	if len(this.ret) == 0 {
		this.ret = [][]int{{val, val}}
		this.arr[val] = true
		return
	}
	if this.arr[val] {
		return
	}
	this.arr[val] = true
	left, right := 0, len(this.ret)
	for left <= right {
		mid := (left + right) / 2
		if val < this.ret[mid][0] {
			if val == this.ret[mid][0]-1 {
				if mid > 0 && this.ret[mid-1][1] == val-1 {
					this.ret[mid-1][1] = this.ret[mid][1]
					this.ret = append(this.ret[:mid], this.ret[mid+1:]...)
					return
				}
				this.ret[mid][0] = val
				return
			}
			if mid > 0 && this.ret[mid-1][1] < val {
				if this.ret[mid-1][1]+1 == val {
					this.ret[mid-1][1] = val
					return
				}
				tmp := append([][]int{{val, val}}, this.ret[mid:]...)
				this.ret = append(this.ret[:mid], tmp...)
				return
			} else if mid == 0 {
				this.ret = append([][]int{{val, val}}, this.ret[mid:]...)
				return
			}
			right = mid - 1
		}
		if val > this.ret[mid][1] {
			if val == this.ret[mid][1]+1 {
				if mid < len(this.ret)-1 && this.ret[mid+1][0] == val+1 {
					this.ret[mid][1] = this.ret[mid+1][1]
					if mid == len(this.ret)-2 {
						this.ret = this.ret[:mid+1]
					} else {
						this.ret = append(this.ret[:mid+1], this.ret[mid+2:]...)
					}
					return
				}
				this.ret[mid][1] = val
				return
			}
			if mid < len(this.ret)-1 && this.ret[mid+1][0] > val {
				if this.ret[mid+1][0]-1 == val {
					this.ret[mid+1][0] = val
					return
				}
				tmp := append([][]int{{val, val}}, this.ret[mid+1:]...)
				this.ret = append(this.ret[:mid+1], tmp...)
				return
			} else if mid == len(this.ret)-1 {
				this.ret = append(this.ret, []int{val, val})
				return
			}
			left = mid + 1
		}
	}
}

func (this *SummaryRanges) GetIntervals() [][]int {
	return this.ret
}

// Your SummaryRanges object will be instantiated and called as such:
/* obj := Constructor();
obj.AddNum(val);
param_2 := obj.GetIntervals();
*/
