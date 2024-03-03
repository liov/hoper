package hard

/*
给你一个整数数组 distance 。

从 X-Y 平面上的点 (0,0) 开始，先向北移动 distance[0] 米，然后向西移动 distance[1] 米，向南移动 distance[2] 米，向东移动 distance[3] 米，持续移动。也就是说，每次移动后你的方位会发生逆时针变化。

判断你所经过的路径是否相交。如果相交，返回 true ；否则，返回 false 。

示例 1：

输入：distance = [2,1,1,2]
输出：true
示例 2：

输入：distance = [1,2,3,4]
输出：false
示例 3：

输入：distance = [1,1,1,1]
输出：true

提示：

1 <= distance.length <= 10^5
1 <= distance[i] <= 10^5

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/self-crossing
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
*/
type point struct {
	x, y int
}

func isSelfCrossing(distance []int) bool {
	m := make(map[point]struct{})
	m[point{0, 0}] = struct{}{}
	last := point{0, 0}
	for i, p := range distance {
		if i%4 == 0 {
			for j := last.y + 1; j <= last.y+p; j++ {
				if _, ok := m[point{last.x, j}]; ok {
					return true
				}
				m[point{last.x, j}] = struct{}{}
			}
			last = point{last.x, last.y + p}
		}
		if i%4 == 1 {
			for j := last.x - 1; j >= last.x-p; j-- {
				if _, ok := m[point{j, last.y}]; ok {
					return true
				}
				m[point{j, last.y}] = struct{}{}
			}
			last = point{last.x - p, last.y}
		}
		if i%4 == 2 {
			for j := last.y - 1; j >= last.y-p; j-- {
				if _, ok := m[point{last.x, j}]; ok {
					return true
				}
				m[point{last.x, j}] = struct{}{}
			}
			last = point{last.x, last.y - p}
		}
		if i%4 == 3 {
			for j := last.x + 1; j <= last.x+p; j++ {
				if _, ok := m[point{j, last.y}]; ok {
					return true
				}
				m[point{j, last.y}] = struct{}{}
			}
			last = point{last.x + p, last.y}
		}
	}
	return false
}

/*
方法一：归纳法（归纳路径交叉的情况）
思路和算法

根据归纳结果，我们发现所有可能的路径交叉的情况只有以下三类：

第 11 类，如上图所示，第 ii 次移动和第 i-3i−3 次移动（包含端点）交叉的情况，例如归纳中的 4-24−2、4-34−3、4-54−5 和 4-64−6。

这种路径交叉需满足以下条件：

第 i-1i−1 次移动距离小于等于第 i-3i−3 次移动距离。
第 ii 次移动距离大于等于第 i-2i−2 次移动距离。

第 22 类，如上图所示，第 55 次移动和第 11 次移动交叉（重叠）的情况，例如归纳中的 5-25−2 和 5-35−3。这类路径交叉的情况实际上是第 33 类路径交叉在边界条件下的一种特殊情况。

这种路径交叉需要满足以下条件：

第 44 次移动距离等于第 22 次移动距离。
第 55 次移动距离大于等于第 33 次移动距离减第 11 次移动距离的差；注意此时第 33 次移动距离一定大于第 11 次移动距离，否则在上一步就已经出现第 11 类路径交叉的情况了。

第 33 类，如上图所示，第 ii 次移动和第 i-5i−5 次移动（包含端点）交叉的情况，例如归纳中的 6-26−2 和 6-36−3。

这种路径交叉需满足以下条件：

第 i-1i−1 次移动距离大于等于第 i-3i−3 次移动距离减第 i-5i−5 次移动距离的差，且小于等于第 i-3i−3 次移动距离；注意此时第 i-3i−3 次移动距离一定大于第 i-5i−5 次移动距离，否则在两步之前就已经出现第 11 类路径交叉的情况了。
第 i-2i−2 次移动距离大于第 i-4i−4 次移动距离；注意此时第 i-2i−2 次移动距离一定不等于第 i-4i−4 次移动距离，否则在上一步就会出现第 33 类路径交叉（或第 22 类路径交叉）的情况了。
第 ii 次移动距离大于等于第 i-2i−2 次移动距离减第 i-4i−4 次移动距离的差。
代码

Python3JavaC#JavaScriptTypeScriptGolangC++

	func isSelfCrossing(distance []int) bool {
	    for i := 3; i < len(distance); i++ {
	        // 第 1 类路径交叉的情况
	        if distance[i] >= distance[i-2] && distance[i-1] <= distance[i-3] {
	            return true
	        }

	        // 第 2 类路径交叉的情况
	        if i == 4 && distance[3] == distance[1] &&
	            distance[4] >= distance[2]-distance[0] {
	            return true
	        }

	        // 第 3 类路径交叉的情况
	        if i >= 5 && distance[i-3]-distance[i-5] <= distance[i-1] &&
	            distance[i-1] <= distance[i-3] &&
	            distance[i] >= distance[i-2]-distance[i-4] &&
	            distance[i-2] > distance[i-4] {
	            return true
	        }
	    }
	    return false
	}

复杂度分析

时间复杂度：O(n)O(n)，其中 nn 为移动次数。

空间复杂度：O(1)O(1)。

方法二：归纳法（归纳路径不交叉时的状态）
思路和算法

根据归纳结果，我们发现当不出现路径交叉时，只可能有以下三种情况：

第 11 种情况：对于每一次移动 ii，第 ii 次移动距离都比第 i-2i−2 次移动距离更长，例如归纳中的 3-33−3、4-94−9、5-85−8 和 6-86−8。
第 22 种情况：对于每一次移动 ii，第 ii 次移动距离都比第 i-2i−2 次移动距离更短，即归纳中的 3-13−1 具有的性质。
第 33 种情况：对于每一次移动 i < ji<j，都满足第 11 种情况；对于每一次移动 i > ji>j，都满足第 22 种情况。
具体地，对于第 33 种情况的第 jj 次移动，有以下三种情况：

第 3.13.1 种情况：第 jj 次移动距离小于第 j-2j−2 次移动距离减去第 j-4j−4 次移动距离的差，例如归纳中的 5-15−1、5-45−4、6-46−4 等。此时，第 j+1j+1 次移动距离需要小于第 j-1j−1 次移动距离才能不出现路径交叉。在边界条件下，这种情况会变为：第 33 次移动距离小于第 11 次移动距离，即归纳中的 3-13−1；第 44 次移动距离小于第 22 次移动距离，即归纳中的 4-14−1、4-44−4 和 4-74−7。
第 3.23.2 种情况：第 jj 次移动距离大于等于第 j-2j−2 次移动距离减去第 j-4j−4 次移动距离的差，且小于等于第 j-2j−2 次移动距离，例如归纳中的 5-55−5、5-65−6、5-75−7 等。此时，第 j+1j+1 次移动距离需要小于第 j-1j−1 次移动距离减去第 j-3j−3 次移动距离的差，才能不出现路径交叉。在边界条件下，这种情况会变为：第 44 次的移动距离等于第 22 次的移动距离且第 33 次的移动距离大于第 11 次的移动距离，即归纳中的 4-84−8。

作者：LeetCode-Solution
链接：https://leetcode-cn.com/problems/self-crossing/solution/lu-jing-jiao-cha-by-leetcode-solution-dekx/
来源：力扣（LeetCode）
著作权归作者所有。商业转载请联系作者获得授权，非商业转载请注明出处。
*/
func isSelfCrossing2(distance []int) bool {
	for i := 3; i < len(distance); i++ {
		// 第 1 类路径交叉的情况
		if distance[i] >= distance[i-2] && distance[i-1] <= distance[i-3] {
			return true
		}

		// 第 2 类路径交叉的情况
		if i == 4 && distance[3] == distance[1] &&
			distance[4] >= distance[2]-distance[0] {
			return true
		}

		// 第 3 类路径交叉的情况
		if i >= 5 && distance[i-3]-distance[i-5] <= distance[i-1] &&
			distance[i-1] <= distance[i-3] &&
			distance[i] >= distance[i-2]-distance[i-4] &&
			distance[i-2] > distance[i-4] {
			return true
		}
	}
	return false
}

func isSelfCrossing3(distance []int) bool {
	n := len(distance)

	// 处理第 1 种情况
	i := 0
	for i < n && (i < 2 || distance[i] > distance[i-2]) {
		i++
	}

	if i == n {
		return false
	}

	// 处理第 j 次移动的情况
	if i == 3 && distance[i] == distance[i-2] ||
		i >= 4 && distance[i] >= distance[i-2]-distance[i-4] {
		distance[i-1] -= distance[i-3]
	}
	i++

	// 处理第 2 种情况
	for i < n && distance[i] < distance[i-2] {
		i++
	}

	return i != n
}
