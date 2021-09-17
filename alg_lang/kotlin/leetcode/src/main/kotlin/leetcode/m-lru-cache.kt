package leetcode

/*
146. LRU 缓存机制
运用你所掌握的数据结构，设计和实现一个  LRU (最近最少使用) 缓存机制 。
实现 LRUCache 类：

LRUCache(int capacity) 以正整数作为容量 capacity 初始化 LRU 缓存
int get(int key) 如果关键字 key 存在于缓存中，则返回关键字的值，否则返回 -1 。
void put(int key, int value) 如果关键字已经存在，则变更其数据值；如果关键字不存在，则插入该组「关键字-值」。当缓存容量达到上限时，它应该在写入新数据之前删除最久未使用的数据值，从而为新的数据值留出空间。


进阶：你是否可以在 O(1) 时间复杂度内完成这两种操作？



示例：

输入
["LRUCache", "put", "put", "get", "put", "get", "put", "get", "get", "get"]
[[2], [1, 1], [2, 2], [1], [3, 3], [2], [4, 4], [1], [3], [4]]
输出
[null, null, null, 1, null, -1, null, -1, 3, 4]

解释
LRUCache lRUCache = new LRUCache(2);
lRUCache.put(1, 1); // 缓存是 {1=1}
lRUCache.put(2, 2); // 缓存是 {1=1, 2=2}
lRUCache.get(1);    // 返回 1
lRUCache.put(3, 3); // 该操作会使得关键字 2 作废，缓存是 {1=1, 3=3}
lRUCache.get(2);    // 返回 -1 (未找到)
lRUCache.put(4, 4); // 该操作会使得关键字 1 作废，缓存是 {4=4, 3=3}
lRUCache.get(1);    // 返回 -1 (未找到)
lRUCache.get(3);    // 返回 3
lRUCache.get(4);    // 返回 4


提示：

1 <= capacity <= 3000
0 <= key <= 10000
0 <= value <= 105
最多调用 2 * 105 次 get 和 put

https://leetcode-cn.com/problems/lru-cache/
 */
class LRUCache(private val capacity: Int) {

  val map = HashMap<Int, LNode>()
  var size = 2
  var head= LNode(0,0)
  var tail= LNode(0,0)
  init{
    head.next = tail
    tail.prev = head
  }

  fun get(key: Int): Int {
    val node = map[key]
    return if (node == null) -1
    else {
      insert(node)
      node.value
    }
  }

  private fun insert(node:LNode){
    if (head != node) {
      head.prev = node
      if (tail == node) tail = node.prev!! else {
        node.prev!!.next = node.next
        node.next!!.prev =node.prev
      }
      node.next = head
      head = node
    }
  }

  fun put(key: Int, value: Int) {
    if(capacity == 1) {
      map.remove(head.key)
      map[key] = head
      head.key = key
      head.value = value
      return
    }
    val node = map[key]
    if (node != null) {
      node.value = value
      insert(node)
      return
    }
    val newNode: LNode
    if (size == capacity) {
      newNode = tail
      map.remove(newNode.key)
      newNode.key = key
      newNode.value = value
      tail = tail.prev!!
      newNode.prev = null
    } else {
      newNode = LNode(key, value)
      size++
    }
    map[key] = newNode
    head.prev = newNode
    newNode.next = head
    head = newNode
  }

}

class LNode(var key: Int, var value: Int) {
  var prev: LNode? = null
  var next: LNode? = null
}
