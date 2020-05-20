package leetcode

import org.junit.jupiter.api.Test
import xyz.hoper.leetcode.Solution
import kotlin.system.measureTimeMillis

class Solution {
  @Test fun lengthOfLongestSubstring() {
    val time1 = measureTimeMillis {
      repeat(100000){
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
}
