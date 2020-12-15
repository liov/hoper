package leetcode

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
