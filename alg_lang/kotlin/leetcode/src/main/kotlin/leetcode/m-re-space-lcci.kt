package leetcode

/**
面试题 17.13. 恢复空格

哦，不！你不小心把一个长篇文章中的空格、标点都删掉了，并且大写也弄成了小写。像句子"I reset the computer. It still didn’t boot!"已经变成了"iresetthecomputeritstilldidntboot"。在处理标点符号和大小写之前，你得先把它断成词语。当然了，你有一本厚厚的词典dictionary，不过，有些词没在词典里。假设文章用sentence表示，设计一个算法，把文章断开，要求未识别的字符最少，返回未识别的字符数。

注意：本题相对原题稍作改动，只需返回未识别的字符数



示例：

输入：
dictionary = ["looked","just","like","her","brother"]
sentence = "jesslookedjustliketimherbrother"
输出： 7
解释： 断句后为"jess looked just like tim her brother"，共7个未识别字符。
提示：

0 <= len(sentence) <= 1000
dictionary中总字符数不超过 150000。
你可以认为dictionary和sentence中只包含小写字母。

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/re-space-lcci
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
//前缀树，hash
//考虑转移方程，每次转移的时候我们考虑第 j(j≤i) 个到第 i 个字符组成的子串sentence[j−1⋯i−1] （注意字符串下标从 0 开始）是否能在词典中找到，如果能找到的话按照定义转移方程即为
//dp[i]=min(dp[i],dp[j−1])
//
//否则没有找到的话我们可以复用 dp[i−1] 的状态再加上当前未被识别的第 i 个字符，因此此时 dp 值为

//dp[i]=dp[i−1]+1

fun respace(dictionary: Array<String>, sentence: String): Int {
  val tree = TrieTree(dictionary)

  val dp = IntArray(sentence.length + 1) { Int.MAX_VALUE }
  dp[0] = 0
  for (i in 1..sentence.length) {
    dp[i] = dp[i - 1] + 1
    var tmp = tree.root
    for (j in i downTo 1) {
      val t = sentence[j - 1] - 'a'
      if (tmp.child == null) break
      else if (tmp.child!![t].isWord) dp[i] = kotlin.math.min(dp[i], dp[j - 1])
      if (dp[i] == 0) break
      tmp = tmp.child!![t]
    }
  }
  return dp.last()
}

class TrieTree {
  val root = TrieTreeNode(false)

  constructor(dictionary: Array<String>) {
    for (i in dictionary.indices) insert(dictionary[i])
  }

  private fun insert(s: String) {
    var tmp = root
    for (i in s.length - 1 downTo 0) tmp = tmp.insert(s[i])
    tmp.isWord = true
  }
}

class TrieTreeNode(var isWord: Boolean) {
  var child: Array<TrieTreeNode>? = null
  fun insert(c: Char): TrieTreeNode {
    if (child == null) child = Array(26) { TrieTreeNode(false) }
    return child!![c - 'a']
  }
}
