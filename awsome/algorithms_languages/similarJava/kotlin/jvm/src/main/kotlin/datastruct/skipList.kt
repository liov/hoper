package datastruct

import kotlin.random.Random

class SkipList<T:Comparable<T>>(init:T){
  private val MAXLEVEL = 16
  private val levelCount = 1
  private val head: Node<T> = Node<T>(init,MAXLEVEL)
  private val random: Random = Random.Default
  var size = 0

  fun add(data:T){

  }

  private fun randomLevel(): Int {
    var level = 0
    for (i in 0 until MAXLEVEL)  if (random.nextInt() % 2 == 1) level++
    return level
  }
  inner class Node<T>(data:T, level:Int){
    var forwards: Array<Node<T>?> = arrayOfNulls(level)
  }

}
