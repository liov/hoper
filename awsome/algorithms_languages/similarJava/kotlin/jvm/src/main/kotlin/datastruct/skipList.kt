package datastruct

import kotlin.random.Random

class SkipList<T : Comparable<T>>(init: T) {
  private val maxLevel = 16
  private var levelCount = 1
  private val head: Node<T> = Node(init, maxLevel)
  private val random: Random = Random.Default
  var size = 0

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
    while (level > maxLevel) level = randomLevel()
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
    val deleteNode: Array<Node<T>?> = arrayOfNulls(maxLevel)
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

  private fun randomLevel(): Int {
    var level = 1
    while (random.nextInt() and 1 == 1) level++
    return level
  }

  fun printAll() {
    var p = head
    while (p.forwards[0] != null) {
      print(p.forwards[0].toString() + " ")
      p = p.forwards[0]!!
    }
    println()
  }

  inner class Node<T>(val data: T, level: Int) {
    var forwards: Array<Node<T>?> = arrayOfNulls(level)
    override fun toString(): String {
      return "Node($data)"
    }
  }
}

fun main(args: Array<String>) {
  val sl = SkipList(-1)
  sl.add(2)
  sl.add(6)
  sl.add(9)
  sl.add(5)
  sl.add(3)
  sl.printAll()
  println(sl[5])
}
