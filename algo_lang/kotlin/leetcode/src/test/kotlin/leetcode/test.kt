package leetcode

import org.junit.jupiter.api.Test
import kotlin.system.measureTimeMillis
import kotlin.time.ExperimentalTime
import kotlin.time.measureTime

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
        LengthOfLongestSubstring.lengthOfLongestSubstringV2("abcabcbb")
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
  fun fourSum() {
    val arr = intArrayOf(1, 0, -1, 0, -2, 2)
    println(fourSum(arr, 0))
  }

  @Test
  fun removeNthFromEnd() {
    val node = ListNode(1).apply {
      next = ListNode(2).apply {
        next = ListNode(3).apply {
          next = ListNode(4).apply { next = ListNode(5) }
        }
      }
    }
    println(removeNthFromEnd(node, 2))
  }


  @Test
  fun combinationSum() {
    val arr = intArrayOf(1, 2)
    println(combinationSum(arr, 4))
  }

  @Test
  fun strStr() {
    println(strStr("mississippi", "pi"))
  }

  @Test
  fun serialize() {
    val node = TreeNode(0).apply {
      left = TreeNode(0).apply {
        left = TreeNode(0)
      }
      right = TreeNode(0).apply {
        right = TreeNode(1).apply {
          right = TreeNode(2)
        }
      }
    }
    println(serialize(deserialize("0,0,0,0,null,null,1,null,null,null,2")))
  }

  @Test
  fun maxScoreSightseeingPair() {
    println(maxScoreSightseeingPair(intArrayOf(8, 1, 5, 2, 6)))
  }

  @Test
  fun recoverFromPreorder() {
    println(serialize(recoverFromPreorder("7-6--2--10---1----7-----4----10---4")))
  }

  @Test
  fun isPalindrome() {
    println(isPalindrome("0P"))
  }

  @Test
  fun rotate() {
    val arr = arrayOf(
      intArrayOf(5, 1, 9, 11),
      intArrayOf(2, 4, 8, 10),
      intArrayOf(13, 3, 6, 7),
      intArrayOf(15, 14, 12, 16)
    )
    rotate(arr)
    println(arr)
  }

  @Test
  fun firstMissingPositive() {
    println(firstMissingPositive(intArrayOf(0, -1, 3, 1)))
  }

  @Test
  fun addBinary() {
    println(addBinary("1010", "1011"))
  }

  @Test
  fun threeSumClosest() {
    println(threeSumClosest(intArrayOf(1, 6, 9, 14, 16, 70), 81))
  }

  @Test
  fun myPow() {
    println(Int.MIN_VALUE)
    println(Int.MAX_VALUE)
    println(myPow(2.00000, -2147483648))
  }

  @Test
  fun searchRange() {
    val arr = searchRange(intArrayOf(2, 2, 2), 2)
    println("${arr[0]},${arr[1]}")
  }

  @Test
  fun minSubArrayLen() {
    println(minSubArrayLen(7, intArrayOf(2, 3, 1, 2, 4, 3)))
  }

  @Test
  fun longestValidParentheses() {
    println(longestValidParentheses("(()("))
  }

  @Test
  fun isValidSudoku() {
    for (i in 0 until 9) {
      println("---$i---")
      for (j in 0 until 9) {
        val x = j / 3 + (i / 3) * 3
        val y = j % 3 + (i % 3) * 3
        println("($x,$y)")
      }
    }
  }
  @Test
  fun findKthLargest(){
    println(findKthLargest(intArrayOf(3,2,1,5,6,4),2))
  }
  @Test
  fun findLength(){
    println(findLength(intArrayOf(3,2,1,5),intArrayOf(1,5)))
  }
  @Test
  fun sortedArrayToBST(){
   sortedArrayToBST(intArrayOf(0,1,2,3,4,5))
  }
  @Test
  fun permute(){
    println(permute(intArrayOf(1,2,3,4)))
  }
  @Test
  fun hasPathSum(){
    println(hasPathSum(deserialize("5,4,8,11,null,13,4,7,2,null,null,null,1"),22))
  }
  @Test
  fun divingBoard(){
    printArr(divingBoard(1,2,5))
  }
  @Test
  fun trieTree(){
    println(respace(arrayOf("looked","just","ju","like","her","brother"),"jesslookedjustliketimherbrother"))
  }
  @Test
  fun intersect(){
    printArr(intersect(intArrayOf(4,9,5),intArrayOf(9,4,8,4)))
  }
  @Test
  fun subsets(){
    println(subsets(intArrayOf(1,2,3,4)))
  }
  @Test
  fun lru(){
    val cache = LRUCache(10)
    cache.put(10, 13)
    cache.put(3, 17)
    cache.put(6, 11)
    cache.put(10, 5)
    cache.get(13)
  }
  @Test
  fun searchInsert(){
    println(searchInsert(intArrayOf(1,3,5,6),5))
  }
  @Test
  fun maxProfit(){
    println(maxProfit(intArrayOf(7,6,4,5,8,10,1,5)))
    println(maxProfitV4(intArrayOf(7,1,5,3,6,6,6)))
  }
  @Test
  fun generateTrees(){
    println(generateTrees(3))
  }
  @Test
  fun findMin(){
    println(findMin(intArrayOf(1,1,3,4,5,6)))
  }
  @Test
  fun minPathSum(){
    println(minPathSum(arrayOf(
      intArrayOf(1,3,1),
      intArrayOf(1,5,1),
      intArrayOf(4,2,1)
    )))
  }
  @Test
  fun isSubsequence(){
    println(isSubsequenceV2("b","c"))
  }
  @Test
  fun integerBreak(){
    println(integerBreak(10))
  }
  @Test
  fun addStrings(){
    println(addStrings("1","1"))
  }
  @ExperimentalTime
  @Test
  fun palindromePairs(){
    val t = arrayOf("abcd","dcba","lls","s","sssll")
    val p = PalindromePairs()
    val t1 = measureTime {
        for(i in 0 until 10000000){
          p.palindromePairs(t)
        }
    }
    val t2 = measureTime {
      for(i in 0 until 10000000){
        palindromePairsV3(t)
      }
    }
    println("$t1,$t2")
  }
}
