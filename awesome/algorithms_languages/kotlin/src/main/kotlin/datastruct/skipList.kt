package datastruct

import kotlin.random.Random

class SkipArrayList<T : Comparable<T>>(init: T) {

  private var levelCount = 1
  private val head: Node<T> = Node(init, MAX_LEVEL)
  var size = 0

  companion object {
    private const val MAX_LEVEL = 16
    private val random: Random = Random.Default
    @JvmStatic private fun randomLevel(): Int {
      var level = 1
      while (random.nextInt() and 1 == 1) level++
      return level
    }
  }

  operator fun get(value: T): Node<T>? {
    var p = head
    for (i in levelCount - 1 downTo 0) {
      while (p.forwards[i] != null && p.forwards[i]!!.data < value) p = p.forwards[i]!!
    }
    return if (p.forwards[0] != null && p.forwards[0]!!.data === value) p.forwards[0] else null
  }

  fun add(data: T) {
    var p = head
    var level = randomLevel()
    while (level > MAX_LEVEL) level = randomLevel()
    val node = Node(data, level)
    val update: Array<Node<T>?> = arrayOfNulls(level)
    for (i in level - 1 downTo 0) {
      while (p.forwards[i] != null && p.forwards[i]!!.data < data) p = p.forwards[i]!!
      update[i] = p
    }
    for (i in 0 until level) {
      node.forwards[i] = update[i]!!.forwards[i]
      update[i]!!.forwards[i] = node
    }
    if (levelCount < level) levelCount = level
    size += 1
  }

  fun remove(value: T) {
    val deleteNode: Array<Node<T>?> = arrayOfNulls(MAX_LEVEL)
    var p = head
    for (i in levelCount - 1 downTo 0) {
      while (p.forwards[i] != null && p.forwards[i]!!.data < value) p = p.forwards[i]!!
      deleteNode[i] = p
    }
    if (p.forwards[0]?.data === value) {
      for (i in levelCount - 1 downTo 0) {
        if (deleteNode[i] != null && deleteNode[i]!!.forwards[i]!!.data === value) {
          deleteNode[i]!!.forwards[i] = deleteNode[i]!!.forwards[i]!!.forwards[i]
        }
      }
    }
  }

  fun printAll() {
    var p = head
    while (p.forwards[0] != null) {
      print(p.forwards[0].toString() + " ")
      p = p.forwards[0]!!
    }
    println()
  }

  class Node<T>(val data: T, level: Int) {
    var forwards: Array<Node<T>?> = arrayOfNulls(level)
    override fun toString(): String {
      return "Node($data)"
    }
  }
}

fun main(args: Array<String>) {
  val srl = SkipArrayList(Int.MIN_VALUE)
  srl.add(2)
  srl.add(6)
  srl.add(9)
  srl.add(5)
  srl.add(3)
  srl.printAll()
  println(srl[5])

  val sll = SkipListList<Int>()
  sll.put(10,5)
  sll.put(2,3)
  sll.put(9,9)
  sll.put(3,6)
  println(sll[9])
  sll.remove(9)
  println(sll[9])
}

class SkipListList<T>{
  private var head = Node<T>(Node.HEAD_KEY, null)
  private var tail = Node<T>(Node.TAIL_KEY, null)
  private var size = 0
  private var listLevel = 0

  init{
    head.next = tail
    tail.pre = head
  }

  companion object {
    private const val PROBABILITY = 0.5
    private val random: Random = Random.Default
  }

  operator fun get(key: Int): Node<T>? {
    val p = findNode(key)
    return if(p.key == key) p else null
  }

  /**
   * put方法有一些需要注意的步骤：
   * 1.如果put的key值在跳跃表中存在，则进行修改操作；
   * 2.如果put的key值在跳跃表中不存在，则需要进行新增节点的操作，并且需要由random随机数决定新加入的节点的高度（最大level）；
   * 3.当新添加的节点高度达到跳跃表的最大level，需要添加一个空白层（除了-oo和+oo没有别的节点）
   *
   * @param k
   * @param v
   */
  fun put(k: Int, v: T) {
    println("添加key:$k")
    var p: Node<T> = findNode(k) //这里不用get是因为下面可能用到这个节点
    println("找到P:$p")
    if (p.key == k) {
      p.value = v
      return
    }
    var q: Node<T> = Node(k, v)
    insertNode(p, q)
    var currentLevel = 0
    while (random.nextDouble() > PROBABILITY) {
      if (currentLevel >= listLevel) addEmptyLevel()
      p = head
      //创建 q的镜像变量（只存储k，不存储v，因为查找的时候会自动找最底层数据）
      val z: Node<T> = Node(k, null)
      insertNode(p, z)
      z.down = q
      q.up = z
      //别忘了把指针移到上一层。
      q = z
      currentLevel++
      println("添加后$this")
    }
    size++
  }


  private fun insertNode(p: Node<T>, q: Node<T>) {
    q.next = p.next
    q.pre = p
    p.next!!.pre = q
    p.next = q
  }

  private fun addEmptyLevel() {
    val p1 = Node<T>(Node.HEAD_KEY, null)
    val p2 = Node<T>(Node.TAIL_KEY, null)
    p1.next = p2
    p1.down = head
    p2.pre = p1
    p2.down = tail
    head.up = p1
    tail.up = p2
    head = p1
    tail = p2
    listLevel++
  }

  //首先查找到包含key值的节点，将节点从链表中移除，接着如果有更高level的节点，则repeat这个操作即可。
  fun remove(k: Int): T? {
    var p = get(k)
    val oldV = p?.value
    var q: Node<T>
    while (p != null) {
      q = p.next!!
      q.pre = p.pre
      p.pre!!.next = q
      p = p.up
    }
    return oldV
  }

  private fun findNode(key: Int): Node<T> {
    var p: Node<T> = head
    while (true) {
      if (p.next != null && p.next!!.key <= key) p = p.next!!
      if (p.down != null)  p = p.down!!
      else if (p.next != null && p.next!!.key > key) break
    }
    return p
  }

  override fun toString():String{
    var p = head
    while (p.down!=null) p = p.down!!
    val sb = StringBuilder()
    while(p.next != null) {
      sb.append(p.toString())
      p = p.next!!
    }
    return sb.toString()
  }

  class Node<T>(val key:Int,var value:T?){
    var pre: Node<T>? = null
    var next:Node<T>? = null
    var up:Node<T>? = null
    var down:Node<T>? = null //上下左右四个节点，pre和up存在的意义在于 "升层"的时候需要查找相邻节点
    companion object {
      const val HEAD_KEY = Int.MIN_VALUE // 负无穷
      const val TAIL_KEY = Int.MAX_VALUE // 正无穷
    }

    override fun toString(): String {
      return "($key,$value)"
    }
  }
}


class SkipList<K:Comparable<K>,T> {
  private var head = SkipListList.Node<T>(SkipListList.Node.HEAD_KEY, null)
  private var tail = SkipListList.Node<T>(SkipListList.Node.TAIL_KEY, null)

  class Node<K,V>(val key:K,var v:V){
    var pre: Node<K,V>? = null
    var next: Node<K,V>? = null
    lateinit var forwards: Array<Node<K,V>?>
  }
}
