package leetcode

/**
34. 在排序数组中查找元素的第一个和最后一个位置

给定一个按照升序排列的整数数组 nums，和一个目标值 target。找出给定目标值在数组中的开始位置和结束位置。

你的算法时间复杂度必须是 O(log n) 级别。

如果数组中不存在目标值，返回 [-1, -1]。

示例 1:

输入: nums = [5,7,7,8,8,10], target = 8
输出: [3,4]
示例 2:

输入: nums = [5,7,7,8,8,10], target = 6
输出: [-1,-1]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/find-first-and-last-position-of-element-in-sorted-array
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun searchRange(nums: IntArray, target: Int): IntArray {
  var left = 0
  var right = nums.size - 1
  val ret = IntArray(2) { -1 }
  if (right<0) return  ret
  var mid:Int
  while (true) {
    mid = (left + right) / 2
    if (left == mid) {
      //压根没找到
      if (nums[left] != target && nums[right] != target) break
      if (ret[1] > mid) {
        right = ret[0]
        left = ret[1]
        //尽量靠左
        ret[0] = if (nums[mid] == target) mid else mid + 1
        continue
      } else {
        //尽量靠右
        ret[1] = if (nums[right] == target) right else left
        if (ret[0] == -1) { //面对[2,2],2的特殊情况
          if(nums[left] == target) ret[0] = left else ret[0] = ret[1]
        }
        return ret
      }
    }
    if (nums[mid] < target) {
      left = mid
      continue
    }
    if (nums[mid] > target) {
      right = mid
      continue
    }
    //第一次找到
    if (ret[0] == -1) {
      //不想新建变量就用这种方式，平时不建议，可读性不好
      //这里必须0记录右索引1记录中位，因为后续遍历是先找的左边索引，会先覆盖左边，但要保留中位做判断
      ret[0] = right
      ret[1] = mid
    }
    //当中点大于第一次找到目标值的中点时说明在找右边索引
    if (ret[1] < mid) left = mid else right = mid
  }
  return ret
}
