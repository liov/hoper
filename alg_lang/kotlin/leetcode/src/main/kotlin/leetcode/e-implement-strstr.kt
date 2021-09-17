package leetcode

/**
28. 实现 strStr()
实现 strStr() 函数。

给你两个字符串 haystack 和 needle ，请你在 haystack 字符串中找出 needle 字符串出现的第一个位置（下标从 0 开始）。如果不存在，则返回  -1 。



说明：

当 needle 是空字符串时，我们应当返回什么值呢？这是一个在面试中很好的问题。

对于本题而言，当 needle 是空字符串时我们应当返回 0 。这与 C 语言的 strstr() 以及 Java 的 indexOf() 定义相符。



示例 1：

输入：haystack = "hello", needle = "ll"
输出：2
示例 2：

输入：haystack = "aaaaa", needle = "bba"
输出：-1
示例 3：

输入：haystack = "", needle = ""
输出：0


提示：

0 <= haystack.length, needle.length <= 5 * 104
haystack 和 needle 仅由小写英文字符组成
 */

fun strStr(haystack: String, needle: String): Int {
  if (needle == "") return 0
  if(haystack.length < needle.length) return -1
  return rabinKarp(haystack,needle)
}

fun KMP(raw: String, pat: String): Int {
  val next = buildNext(pat)
  val n = raw.length
  var i = 0
  val m = pat.length
  var j = 0
  while (j < m && i < n) {
    if (0 > j || raw[i] == pat[j]) {
      i++
      j++
    } else j = next[j]
  }
  return if (j == m) i - j else -1
}

fun buildNext(pat: String): IntArray {
  val m = pat.length
  var j = 0
  val next = IntArray(m)
  var t = -1
  next[0] = t
  while (j < m - 1) {
    if (0 > t || pat[j] == pat[t]) {
      j++
      t++
      next[j] = if (pat[j] != pat[t]) t else next[t]
    } else t = next[t]
  }
  return next
}

fun Sunday(raw: String, pat: String): Int {
  val n = raw.length
  var i = 0
  val m = pat.length
  var j = 0
  val offset = HashMap<Char, Int>()
  for (c in pat.indices) {
    offset[pat[c]] = m - c
  }
  while (i <= n - m) {
    j = 0
    while (raw[i + j] == pat[j]) {
      j += 1
      if (j == m) return i
    }
    if (i + m == n) return -1
    i += if (raw[i + m] in offset)
      offset[raw[i + m]]!!
    else m + 1
  }
  return -1
}

fun preBmBc(x: String, bmBc: IntArray) {
  // 计算字符串中每个字符距离字符串尾部的长度
  for (i in x.indices) bmBc[x[i].toInt()] = x.length - i - 1
}

// 计算以i为边界，与模式串后缀匹配的最大长度（区间的概念）
fun suffix(x: String, suff: IntArray) {
  val len = x.length
  var q: Int
  for (i in len - 2 downTo 0) {
    q = i
    while (q >= 0 && x[q] == x[len - 1 - i + q]) {
      --q;
    }
    suff[i] = i - q;
  }
}

// 好后缀算法的预处理
/*
 有三种情况
 1.模式串中有子串匹配上好后缀
 2.模式串中没有子串匹配上好后缀，但找到一个最大前缀
 3.模式串中没有子串匹配上好后缀，但找不到一个最大前缀


 3种情况获得的bmGs[i]值比较

 3 > 2 > 1

 为了保证其值越来越小

 所以按顺序处理3->2->1情况
 */
fun preBmGs(s: String, bmGs: IntArray) {

  val len = s.length
  val suff = IntArray(len)

  suffix(s, suff)

  //全部初始为自己的长度，处理第三种情况
/*
  for (i in s.indices) {
    bmGs[i] = len
  }
*/

  // 处理第二种情况
  for (i in len - 1 downTo 0) {
    if (suff[i] == i + 1) { // 找到合适位置
      for (j in s.indices) if (bmGs[j] == len) bmGs[j] = len - 1 - i // 保证每个位置至多只能被修改一次
    }
  }

  // 处理第一种情况，顺序是从前到后
  for (i in s.indices)  bmGs[len - 1 - suff[i]] = len - 1 - i

}

fun BM(raw: String, pat: String): Int {

  val n = raw.length
  val m = pat.length

  val bmGs = IntArray(m) { m }
  // 全部更新为自己的长度
  val bmBc = IntArray(256) { m }

  // 处理好后缀算法
  preBmGs(pat, bmGs)
  // 处理坏字符算法
  preBmBc(pat, bmBc)

  var j = 0

  while (j <= n - m) {
    // 模式串向左边移动
    var i = m - 1
    while (i >= 0 && pat[i] == raw[i + j]) i--
    // 给定字符串向右边移动
    if (i < 0)  return j// 移动到模式串的下一个位置
   else {
      // 取移动位数的最大值向右移动，前者好后缀，向右移动，后者坏字符，向左移动
      j += kotlin.math.max(bmGs[i], bmBc[raw[i + j].toInt()] - m + 1 + i);
    }
  }
  return -1
}

//jvm没有无符号整型
//go的子串算法
@ExperimentalUnsignedTypes
fun RabinKarp(s: String, substr: String): Int {
  // Rabin-Karp search
  val (hashss, pow) = hashStr(substr)
  val n = substr.length
  var h = 0u
  for (i in substr.indices) h = h * 16777619u + s[i].toInt().toUInt()
  if (h == hashss && s.substring(0, n) == substr)  return 0

  for (i in n until s.length) {
    h = h*16777619u + s[i].toInt().toUInt() - pow * s[i - n].toInt().toUInt()
    if (h == hashss && s.substring(i - n, i) == substr) return i - n
  }
  return -1
}

@ExperimentalUnsignedTypes
fun hashStr(sep: String): StrHash {
  var hash = 0u
  for (i in sep.indices)  hash = hash * 16777619u + sep[i].toInt().toUInt()
  var pow = 1u
  var sq = 16777619u
  var i = sep.length
  while (i > 0) {
    i = i shr 1
    if (i and 1 != 0)  pow *= sq
    sq *= sq
  }
  return StrHash(hash, pow)
}
@ExperimentalUnsignedTypes
data class StrHash(val hash: UInt, val pow: UInt)

@ExperimentalUnsignedTypes
fun rabinKarp(str:String, pattern:String) :Int{
  val n = str.length
  val m = pattern.length

  //哈希时需要用到进制计算，这里只涉及26个字母所以使用26进制
  val d = 26u
  //防止hash之后的值超出int范围，对最后的hash值取模
  //q取随机素数，满足q*d < INT_MAX即可
  val q = 16777619u

  //str子串的hash值
  var strCode = 0u
  //pattern的hash值
  var patternCode = 0u
  //d的size2-1次幂，hash计算时，公式中会用到
  var h = 1u

  //计算sCode、pCode、h
  for (i in 0 until m) {
    patternCode = (d*patternCode + (pattern[i]-'a').toUInt()) % q
    //计算str第一个子串的hash
    strCode = (d*strCode + (str[i]-'a').toUInt()) % q
    h = (h*d) % q
  }
  //最大需要匹配的次数
  //字符串开始匹配，对patternCode和strCode开始比较，并更新strCode的值
  for (i in 0 until n - m + 1) {
    if(strCode == patternCode && ensureMatching(i, str, pattern)) return i
    if(i == n-m) break
    //更新strCode的值，即计算str[i+1,i+m-1]子串的hashCode
    strCode = (strCode*d - h*(str[i]-'a').toUInt() + (str[i+m] - 'a').toUInt())%q
  }
  return -1
}

/**
 * hash值一样并不能完全确保字符串一致，所以还需要进一步确认
 * @param i hash值相同时字符串比对的位置
 * @param pattern 模式串
 * @return
 */
fun ensureMatching(i:Int, str:String,pattern:String) :Boolean {
  val strSub = str.substring(i, i+pattern.length)
  return strSub == pattern;
}
