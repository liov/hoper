package leetcode

import java.util.*
import kotlin.collections.ArrayList

/**
336. 回文对

给定一组唯一的单词， 找出所有不同的索引对(i, j)，使得列表中的两个单词，words[i] + words[j]，可拼接成回文串。

示例 1:

输入: ["abcd","dcba","lls","s","sssll"]
输出: [[0,1],[1,0],[3,2],[2,4]]
解释: 可拼接成的回文串为 ["dcbaabcd","abcddcba","slls","llssssll"]
示例 2:

输入: ["bat","tab","cat"]
输出: [[0,1],[1,0]]
解释: 可拼接成的回文串为 ["battab","tabbat"]

来源：力扣（LeetCode）
链接：https://leetcode-cn.com/problems/palindrome-pairs
著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
 */
fun palindromePairs(words: Array<String>): List<List<Int>> {
  val ret = ArrayList<List<Int>>()
  for (i in words.indices) {
    for (j in words.indices) {
      if (i == j) continue
      if (check(words[i], words[j])) ret.add(listOf(i, j))
    }
  }
  return ret
}

fun check(words1: String, words2: String): Boolean {
  val words = words1 + words2
  if (words.length == 1) return true
  for (i in 0 until words.length / 2) {
    if (words[i] != words[words.length - 1 - i]) return false
  }
  return true
}

fun palindromePairsV2(words: Array<String>): List<List<Int>> {
  val trie1 = Trie()
  val trie2 = Trie()
  for (i in words.indices) {
    trie1.insert(words[i], i)
    trie2.insert(words[i].reversed(), i)
  }
  val ret = ArrayList<List<Int>>()
  for (i in words.indices) {
    val rec = manacher(words[i])

    val id1 = trie2.query(words[i])
    words[i] = StringBuffer(words[i]).reverse().toString()
    val id2 = trie1.query(words[i])

    val m = words[i].length

    val allId = id1[m]
    if (allId != -1 && allId != i) ret.add(listOf(i, allId))
    for (j in 0 until m) {
      if (rec[j][0] != 0) {
        val leftId = id2[m - j - 1]
        if (leftId != -1 && leftId != i) ret.add(listOf(leftId, i))
      }
      if (rec[j][1] != 0) {
        val rightId = id1[j]
        if (rightId != -1 && rightId != i) ret.add(listOf(i, rightId))
      }
    }
  }
  return ret
}

/**
 * 回文分为奇回文（ababa）和偶回文（abba），这里比较难以处理，我们使用一个小(sao)技(cao)巧(zuo)（很重要）。我们将字符串首尾和每个字符间插入一个字符（注意：这个自符在串中并未出现）例如：'#'.

栗子！栗子！s='abbadcacda'先转化成s_new='$#a#b#b#a#d#c#a#c#d#a#\0'（'$'与'\0'，是边界，下面的代码中可以看到）

这样原串中的偶回文（abba）与奇回文（adcacda），变成了（#a#d#d#a#）与（#a#d#c#a#c#d#a#）两个奇回文。
那，p[i]该如何求呢？很显然，p[i]-1正好就是原字符中的最长回文串长度了。
 */
fun manacher(s: String): Array<IntArray> {
  val n = s.length
  val tmp = StringBuffer("#")
  if (n > 0) tmp.append(s[0])
  for (i in 1 until n) {
    tmp.append('*')
    tmp.append(s[i])
  }
  tmp.append('!')
  val m = n * 2
  val len = IntArray(m)
  val ret = Array(n) { IntArray(2) }
  var p = 0
  var maxn = -1
  for (i in 1 until m) {
    len[i] = if (maxn >= i) kotlin.math.min(len[2 * p - i], maxn - i) else 0
    while (tmp[i - len[i] - 1] == tmp[i + len[i] + 1]) {
      len[i]++
    }
    if (i + len[i] > maxn) {
      p = i
      maxn = i + len[i]
    }
    if (i - len[i] == 1) {
      ret[(i + len[i]) / 2][0] = 1
    }
    if (i + len[i] == m - 1) {
      ret[(i - len[i]) / 2][1] = 1
    }
  }
  return ret
}

class Trie {
  class Node {
    val ch = IntArray(26)
    var flag = -1
  }

  private val tree = mutableListOf(Node())

  fun insert(s: String, id: Int) {
    var add = 0
    for (i in s.indices) {
      val x = s[i] - 'a'
      if (tree[add].ch[x] == 0) {
        tree.add(Node())
        tree[add].ch[x] = tree.size - 1
      }
      add = tree[add].ch[x]
    }
    tree[add].flag = id
  }

  fun query(s: String): IntArray {
    val len = s.length
    var add = 0
    val ret = IntArray(s.length + 1) { -1 }
    for (i in s.indices) {
      ret[i] = tree[add].flag
      val x = s[i] - 'a'
      if (tree[add].ch[x] == 0) return ret
      add = tree[add].ch[x]
    }
    ret[len] = tree[add].flag
    return ret
  }
}

fun palindromePairsV3(words: Array<String>): List<List<Int>>? {
  //构建trie
  val root = buildTire(words)
  val res: MutableList<List<Int>> = java.util.ArrayList()
  //对于每个单词,在trie中搜索
  for (i in words.indices) {
    search(words[i], i, root, res)
  }
  return res
}

fun search(word: String, i: Int, node: TrieNode, res: MutableList<List<Int>>) {
  var node: TrieNode? = node
  val k = word.length
  var j = 0
  while (j < k) {

    //循环中在trie中发现的单词是比当前word长度要短的
    val c = word[j]
    if (node!!.index != -1 && isPalindrome(word, j, k - 1)) {
      res.add(listOf(i, node.index))
    }
    //所有可能被排除，提前退出
    if (node.children[c - 'a'] == null) return
    node = node.children[c - 'a']
    j++
  }
  //长度等于当前搜索的word的单词
  if (node!!.index != -1 && node.index != i) {
    res.add(listOf(i, node.index))
  }
  //长度长于当前搜索的word的单词
  for (rest in node.belowIsPalindrome) {
    res.add(listOf(i, rest))
  }
}

fun buildTire(words: Array<String>): TrieNode {
  val root = TrieNode()
  for (i in words.indices) {
    addWord(root, words[i], i)
  }
  return root
}

fun addWord(root: TrieNode, word: String, i: Int) {
  var root: TrieNode? = root
  for (j in word.length - 1 downTo 0) {
    if (isPalindrome(word, 0, j)) {
      root!!.belowIsPalindrome.add(i)
    }
    val c = word[j]
    if (root!!.children[c - 'a'] == null) {
      root.children[c - 'a'] = TrieNode()
    }
    root = root.children[c - 'a']
  }
  root!!.index = i
}

fun isPalindrome(word: String, i: Int, j: Int): Boolean {
  var i = i
  var j = j
  if (word.length <= 1) return true
  while (i < j) {
    if (word[i++] != word[j--]) return false
  }
  return true
}

class TrieNode {
  var index = -1
  var belowIsPalindrome =ArrayList<Int>()
  var children: Array<TrieNode?> = arrayOfNulls(26)
}
