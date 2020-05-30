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
    val nums = intArrayOf(1, 1, -2)
    println(threeSum(nums))
  }

  @Test
  fun mergeTwoLists() {
    val node1 = ListNode(1).apply { next = ListNode(2).apply { next = ListNode(4) } }
    val node2 = ListNode(1).apply { next = ListNode(3).apply { next = ListNode(4) } }
    println(mergeTwoLists(node1, node2))
  }

  @Test
  fun isValid() {
    println(isValidV2("{[]}"))
  }

  @Test
  fun generateParenthesis() {
    println(generateParenthesis(3))
  }

  @Test
  fun binaryTree() {
    val arr = intArrayOf(5, 3, 12, 36, 728, 333, 128)
    val bt = BinaryTree<Int>()
    for (i in arr.indices) {
      bt.insert(arr[i])
    }
    bt.prevRecursive()
    println("\n")
    print("${bt.midIterator()}\n")
    print("${bt.prevIterator()}\n")
    bt.subRecursive()
    bt.sequence().forEach { println(it) }
  }

  @Test
  fun mergeKLists() {
    val node1 = ListNode(1).apply { next = ListNode(4).apply { next = ListNode(5) } }
    val node2 = ListNode(1).apply { next = ListNode(3).apply { next = ListNode(4) } }
    val node3 = ListNode(2).apply { next = ListNode(6) }
    val list: Array<ListNode?> = arrayOf(node1, node2, node3)
    println(mergeKListsV2(list))
  }

  @Test
  fun fourSum(){
    val arr = intArrayOf(1, 0, -1, 0, -2, 2)
    println(fourSum(arr,0))
  }

  @Test
   fun removeNthFromEnd(){
    val node = ListNode(1).apply { next = ListNode(2).
    apply { next = ListNode(3).
    apply { next = ListNode(4).
    apply { next = ListNode(5)}}} }
    println(removeNthFromEnd(node,2))
  }


  @Test
  fun combinationSum(){
    val arr = intArrayOf(1,2)
    println(combinationSum(arr,4))
  }
}
