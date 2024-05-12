package leetcode


/**
3. 无重复字符的最长子串

给定一个字符串，请你找出其中不含有重复字符的 最长子串 的长度。

示例 1:

输入: "abcabcbb"
输出: 3
解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
示例 2:

输入: "bbbbb"
输出: 1
解释: 因为无重复字符的最长子串是 "b"，所以其长度为 1。
示例 3:

输入: "pwwkew"
输出: 3
解释: 因为无重复字符的最长子串是 "wke"，所以其长度为 3。
     请注意，你的答案必须是 子串 的长度，"pwke" 是一个子序列，不是子串。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/longest-substring-without-repeating-characters
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//通过	212 ms	33.9 MB	Kotlin
fun lengthOfLongestSubstring(s: String): Int {
  if (s.length <= 1) {
    return s.length
  }

  var maxLength = 0
  var currentLength = 0
  //因为题目都是小写英文字母,数组赋值居然是逐个赋值，又是开销IntArray(26){-1}
  //for 最浪费时间，不建议（其实memset内部也是用循环实现的，只不过memset经过了严格优化，所以性能更高
  val intArray = IntArray(128)
  var idx: Int
  for (i in s.indices) {
    idx = s[i].toInt()
    if (intArray[idx] != 0 && intArray[idx] >= i - currentLength) {
      if (maxLength < currentLength) {
        maxLength = currentLength
      }
      //当前长度等于上次的长度减去（该字符上次出现位置减去新子串第一个元素的索引） 例如bbtablud遍历到第三个b实际减0
      currentLength -= intArray[idx] - (i - currentLength) - 1 //为了避免数组intArray初始化，采取+-1的方式还原真实索引
    } else {
      currentLength++
    }
    intArray[idx] = i + 1
  }
  return if (maxLength > currentLength) maxLength else currentLength
}

/**
 *  这是真的吗
 * 	通过	3 ms	39.9 MB	Java
 *	通过	296 ms	34.1 MB	Kotlin
 */
//优化版 204 ms	33.8 MB	Kotlin
fun lengthOfLongestSubstringV2(s: String): Int {
  if (s.length <= 1) {
    return s.length
  }
  var maxLength = 0
  var left = 0
  val intArray = IntArray(128)
  var idx: Int
  for (i in s.indices) {
    idx = s[i].toInt()
    if (intArray[idx] - 1 > left) left = intArray[idx] - 1
    intArray[idx] = i + 1
    if (i - left > maxLength) maxLength = i - left
  }
  return maxLength
}
