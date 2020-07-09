package leetcode

/**
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
fun respace(dictionary: Array<String>, sentence: String): Int {
  val tree = TrieTree(dictionary)
  var tmp = tree.root
  var nextCount = 0
  var count = sentence.length
  var tmpCount = 0
  var i = 0
  while (i < sentence.length) {
    if (tmp.child == null) {
      tmpCount = 0
      tmp = tree.root
      continue
    }
    if (tmp.child!![26].exists) {
      nextCount = tmpCount
    }
    tmpCount++
    if (tmp.child!![sentence[i] - 'a'].exists) tmp = tmp.child!![sentence[i] - 'a']
    else {
      if (nextCount != 0) {
        count -= nextCount
        nextCount = 0
      }
      if (tmp != tree.root) {
        tmp = tree.root
        i--
      }
      tmpCount = 0
    }
    i++
  }
  if (tmp.child!![26].exists) {
    count -= tmpCount
  }
  return count
}

class TrieTree {
  val root = TrieTreeNode('\u0000', false)

  constructor(dictionary: Array<String>) {
    for (i in dictionary.indices) insert(dictionary[i])
  }

  fun insert(s: String) {
    var tmp = root
    for (i in s.indices) {
      tmp = tmp.insert(s[i])
    }
    tmp.insert('a' + 26)
  }
}

class TrieTreeNode(val c: Char, var exists: Boolean) {
  var child: Array<TrieTreeNode>? = null
  fun insert(c: Char): TrieTreeNode {
    if (child == null) child = Array(27) { i -> TrieTreeNode('a' + i, false) }
    if (!child!![c - 'a'].exists) child!![c - 'a'].exists = true
    return child!![c - 'a']
  }
}

internal class Trie {
  var root: TrieNode

  // 将单词倒序插入字典树
  fun insert(word: String) {
    var cur: TrieNode? = root
    for (i in word.length - 1 downTo 0) {
      val c = word[i] - 'a'
      if (cur!!.children[c] == null) {
        cur.children[c] = TrieNode()
      }
      cur = cur.children[c]
    }
    cur!!.isWord = true
  }

  // 找到 sentence 中以 endPos 为结尾的单词，返回这些单词的开头下标。
  fun search(sentence: String, endPos: Int): List<Int> {
    val indices: ArrayList<Int> = ArrayList()
    var cur: TrieNode? = root
    for (i in endPos downTo 0) {
      val c = sentence[i] - 'a'
      if (cur!!.children[c] == null) break
      cur = cur.children[c]
      if (cur!!.isWord) indices.add(i)
    }
    return indices
  }

  init {
    root = TrieNode()
  }
}

internal class TrieNode {
  var isWord = false
  var children = arrayOfNulls<TrieNode>(26)
}
