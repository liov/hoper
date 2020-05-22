package leetcode

import org.junit.jupiter.api.Test
import xyz.hoper.leetcode.Solution
import kotlin.system.measureTimeMillis

class Solution {
  @Test
  fun lengthOfLongestSubstring() {
    val time1 = measureTimeMillis {
      repeat(100000) {
        lengthOfLongestSubstring("abcabcbb")
      }

    }
    val time2 = measureTimeMillis {
      repeat(100000) {
        Solution.lengthOfLongestSubstringV2("abcabcbb")
      }
    }
    //kotlin还快一点
    println("time:$time1")
    println("time:$time2")
  }

  @Test
  fun reverse() {
    println(reverse(1534236469))
  }

  @Test
  fun findMedianSortedArrays() {
    val nums1 = intArrayOf(3)
    val nums2 = intArrayOf(-2, -1)
    println(findMedianSortedArrays(nums1, nums2))
  }

  @Test
  fun threeSum() {
    val nums = intArrayOf(1,1,-2)
    println(threeSum(nums))
  }

  @Test
  fun mergeTwoLists(){
    val node1 = ListNode(1).apply { next = ListNode(2).apply { next = ListNode(4)}}
    val node2 = ListNode(1).apply { next = ListNode(3).apply { next = ListNode(4)}}
    println(mergeTwoLists(node1,node2))
  }

  @Test
  fun isValid(){
    println(isValidV2("{[]}"))
  }
}
