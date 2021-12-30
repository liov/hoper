package leetcode



/**
718. 最长重复子数组

给两个整数数组A和B，返回两个数组中公共的、长度最长的子数组的长度。

示例 1:

输入:
A: [1,2,3,2,1]
B: [3,2,1,4,7]
输出: 3
解释:
长度最长的公共子数组是 [3, 2, 1]。
说明:

1 <= len(A), len(B) <= 1000
0 <= A[i], B[i] < 100

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/maximum-length-of-repeated-subarray
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */

const val mod: Int = 1000000009
const val base = 113

fun findLength(A: IntArray, B: IntArray): Int {
  var left = 1
  var right = kotlin.math.min(A.size, B.size) + 1
  while (left < right) {
    val mid = left + right shr 1
    if (check(A, B, mid))  left = mid + 1 else right = mid
  }
  return left - 1
}
//字符串匹配，Rabin-Karp
//其中，base代表的是S的定义域大小，比如说如果S全是英文字母，那么的值为26，因为英文字母就只有26个。然后这个函数是一个映射函数，映射S的定义域中的每一个字符到数字的函数。
//如果采用O(m)的算法计算长度为m的字符串子串的哈希值的话，那复杂度还是O(nm)。这是就要使用一个滚动哈希的优化技巧。
//选取两个合适的互素常数b和h(l<b<h)，假设字符串C=c1c2...cm，定义哈希函数：
//H(C)=(c1*b^(m-1)+c2*b^(m-2))+c3*b^(m-3)+...+cm*b^0)%h
//其中b是基数，相当于把字符串看作b进制数。这样，字符串S=s1s2...sn从位置k+1开始长度为m的字符串子串S[k+1...k+m]的哈希值，就可以利用从位置k开始的字符串子串S[k...k+m-1]的哈希值直接进行计算：
//H(S[k+1...k+m])=(H(s[k:k+m-1])*b-sk*b^m+s(k+m))%h
//二分+hash，逐步缩小
//hash(S)= i=0∑∣S∣−1 base^∣S∣−(i+1)×S[i]
//hash(S[1:len+1])=(hash(S[0:len])−base^(len−1)×S[0])×base+S[len+1]
fun check(A: IntArray, B: IntArray, len: Int): Boolean {
  var hashA: Long = 0
  //0到len-1的hash
  for (i in 0 until len) hashA = (hashA * base + A[i]) % mod
  val bucketA: MutableSet<Long> = HashSet()
  bucketA.add(hashA)
  val mult = qPow(base.toLong(), len.toLong())
  //1到len，2到len+1...
  for (i in len until A.size) {
    hashA = (hashA* base - A[i - len] * mult + A[i]) % mod
    if (hashA < 0) hashA+=mod 
    bucketA.add(hashA)
  }
  var hashB: Long = 0
  for (i in 0 until len) hashB = (hashB * base + B[i]) % mod
  if (bucketA.contains(hashB)) return true
  for (i in len until B.size) {
    hashB = (hashB* base - B[i - len] * mult + B[i]) % mod
    if (hashB < 0) hashA+=mod
    if (bucketA.contains(hashB)) return true
  }
  return false
}

// 使用快速幂计算 x^n % mod 的值
fun qPow(x: Long, n: Long): Long {
  var x = x
  var n = n
  var ret: Long = 1
  while (n != 0L) {
    if (n and 1 != 0L)  ret = ret * x % mod
    x = x * x % mod
    n = n shr 1
  }
  return ret
}
