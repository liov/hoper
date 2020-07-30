package leetcode

import xyz.hoper.test.oop.swap
import java.util.*
import kotlin.system.measureTimeMillis


//用树还真不好实现
class MaxHeap<T : Comparable<T>> {
  //不可以创建泛型数组，又是JVM语言的一个黑点
  var heap: MutableList<T>
  var size: Int = 0
  private var cap: Int = 0

  constructor(init: T) : this(10, init)

  constructor(cap: Int, init: T) {
    this.cap = cap
    this.heap = MutableList(cap) { init }
  }

  constructor(list: MutableList<T>) {
    this.size = list.size
    this.heap = list
    for (i in size / 2 - 1 downTo 0) {
      adjustHeap(list, size, i)
    }
  }

  private fun adjustHeap(list: MutableList<T>, size: Int, i: Int) {
    val leftChild = 2 * i + 1    //左子节点索引
    val rightChild = 2 * i + 2   //右子节点索引
    var max = i             //选出当前结点与其左右孩子三者之中的最大值
    if (leftChild < size && list[leftChild] > list[max]) max = leftChild

    if (rightChild < size && list[rightChild] > list[max]) max = rightChild

    if (max != i) {
      list.swap(i, max)//将最大值节点与父节点互换
      adjustHeap(list, size, max) //递归调用，继续从当前节点向下进行堆调整
    }
  }
}


private fun parent(pos: Int): Int {
  return (pos - 1) / 2
}

private fun leftChild(pos: Int): Int {
  return pos * 2 + 1
}

private fun rightChild(pos: Int): Int {
  return pos * 2 + 2
}


/***
 * 按树结构打印
 */
fun printHeap(heap: IntArray, size: Int) {
  println(heap.joinToString(prefix="[",postfix="]",limit=size-1,truncated=heap[size-1].toString()))
}

fun topN(less: (Int, Int) -> Boolean) {
  val len = 20
  val topN = 7
  val heap = IntArray(len)
  val random = Random()
  for (i in 0 until len)  heap[i] = random.nextInt(10000)

  //创建topN堆
  createHeap(heap, topN, less)
  printArr(heap)
  //遍历数组，并维护topN堆
  findTopNInPlace(heap, topN, less)
  printHeap(heap, topN)
  printArr(heap)
  printArrSorted(heap)
}

/***
 * 前N个最大值
 * 最小堆
 */

fun topMaxN() = topN() { i, j -> i < j }


/***
 * 前N个最小值
 * 最大堆
 */

fun topMinN() = topN() { i, j -> i > j }

fun<T:Comparable<T>> createHeapUp(heap: MutableList<T>) {
  for (i in 1 until heap.size) adjustUp(heap, i)
}

fun<T:Comparable<T>> createHeapDown(heap: MutableList<T>) {
  for (i in heap.size/2-1 downTo 0 ) adjustDown(heap, i)
}

private fun<T:Comparable<T>> adjustUp(heap: MutableList<T>, pos: Int) {
  var pos = pos
  //比父节点大就上移
  while (parent(pos) >= 0 && heap[pos] > heap[parent(pos)]) {
    val parent = parent(pos)
    heap.swap(parent, pos)
    pos = parent
  }
}

private fun<T:Comparable<T>> adjustDown(heap: MutableList<T>, pos: Int) {
  var pos = pos
  while (leftChild(pos) < heap.size) {
    var child = leftChild(pos)
    //判断左右孩子的大小，child代表较小的孩子
    if (child + 1 < heap.size && heap[child + 1] > heap[child]) child++
    //新节点比较小的孩子都小，说明找到对应位置，直接跳出循环
    if (heap[pos] > heap[child]) break
    heap.swap(pos, child)
    pos = child
  }
}

fun createHeap(heap: IntArray, size: Int, less: (Int, Int) -> Boolean) {
  for (i in 1 until size) adjustUp(heap, i, less)
}

private fun adjustUp(heap: IntArray, pos: Int, less: (Int, Int) -> Boolean) {
  var pos = pos
  //比父节点大就上移
  while (parent(pos) >= 0 && less(heap[pos], heap[parent(pos)])) {
    val parent = parent(pos)
    heap.swap(parent, pos)
    pos = parent
  }
}

/***
 * 遍历数据组，并维护最大堆
 */
private fun findTopNInPlace(heap: IntArray, topN: Int, less: (Int, Int) -> Boolean) {
  for (i in topN until heap.size) adjustDownTopInPlace(heap, i, less)
}

private fun adjustDown(heap: IntArray, i: Int, less: (Int, Int) -> Boolean) {
  var pos = i
  while (leftChild(pos) < heap.size) {
    var child = leftChild(pos)
    //判断左右孩子的大小，child代表较小的孩子
    if (child + 1 < heap.size && less(heap[child + 1], heap[child])) child++
    //新节点比较小的孩子都小，说明找到对应位置，直接跳出循环
    if (less(heap[pos], heap[child])) break
    heap.swap(pos, child)
    pos = child
  }
}

private fun adjustDownTop(heap: IntArray, v: Int, less: (Int, Int) -> Boolean) {
  //比topN中最小的还要小直接返回
  if (less(v, heap[0])) return
  heap[0] = v
  adjustDown(heap, 0, less)
}


/***
 * 向下调整新加入节点位置，并维护最大堆
 */
private fun adjustDownTopInPlace(heap: IntArray, pos: Int, less: (Int, Int) -> Boolean) {
  //比topN中最大的还要大直接返回
  if (less(heap[pos], heap[0])) return
  heap.swap(0, pos)
  adjustDown(heap, 0, less)
}

/***
 * 向下调整新加入节点位置，并维护最小堆
 */
private fun adjustDownTopMaxNInPlace(heap: IntArray, pos: Int) = adjustDownTopInPlace(heap, pos) { i, j -> i < j }

/***
 * 向下调整新加入节点位置，并维护最大堆
 */
private fun adjustDownTopMinNInPlace(heap: IntArray, topN: Int, pos: Int) =  adjustDownTopInPlace(heap, pos) { i, j -> i > j }

/***
 * 遍历数据组，并维护最小堆
 */
fun findTopMaxNInPlace(heap: IntArray, topN: Int) = findTopNInPlace(heap, topN) { i, j -> i < j }

/***
 * 遍历数据组，并维护最大堆
 */
private fun findTopMinNInPlace(heap: IntArray, topN: Int) = findTopNInPlace(heap, topN) { i, j -> i > j }

private fun adjustDownTopMaxN(heap: IntArray, v: Int) = adjustDownTop(heap, v) { i, j -> i < j }


//原来一直疑惑topN为什么用小顶堆，我还想要用最大堆，然后插进来一个元素就剔除最小的那个，原来最小堆直接剔除根节点就行了...
/*首先，我们需要构建一个大小为N的小顶堆，小顶堆的性质如下：
每一个父节点的值都小于左右孩子节点，然后依次从文件中读取10亿个整数，如果元素比堆顶小，则跳过不进行任何操作，
如果比堆顶大，则把堆顶元素替换掉，并重新构建小顶堆。当10亿个整数遍历完成后，堆内元素就是TopN的结果。*/
fun main(args: Array<String>) {
  topMaxN()
  //topMinN()
  val toN = TopN(MutableList(9) { it * 2 + 1 })
  val random = Random()
  val list = ArrayList<Int>()
  var v: Int
  for (i in 0 until 20) {
    v = random.nextInt(1000)
    toN.insert(v)
    list.add(v)
  }
  val mh = MaxHeap(list.clone() as MutableList<Int>)
  println(mh.heap)
 val list1 = list.clone() as MutableList<Int>
  createHeapUp(list1)
  println(list1)
  val list2 = list.clone() as MutableList<Int>
  createHeapDown(list2)
  println(list2)
  list.sort()
  for (i in 11..19)  print("${list[i]}, ")
  println()
  toN.list.sort()
  println(toN.list)

  val arr1 = IntArray(100) { random.nextInt(100000000) }
  val arr2 = arr1.clone()
  val heap1 = IntArray(9)
  val heap2 = TopN(MutableList(9) { 0 })
  val time1 = measureTimeMillis {
    for ((_, value) in arr1.withIndex()) adjustDownTopMaxN(heap1, value)
  }

  printArr(heap1)
  //好吧，费了两天时间写的环形优先级队列性能更差，10倍左右
  //2020/07/16 今天怎么测的性能更好了，两倍，难道lambda的开销？
  val time2 = measureTimeMillis {
    for ((_, value) in arr2.withIndex())  heap2.insert(value)
  }
  println(heap2.list)
  println("$time1,$time2")
}

//昨晚想的链表加数组，好像不需要链表
class TopN<T : Comparable<T>>(val list: MutableList<T>) {
  var minIdx = 0

  init {
    //list.sort()
  }
  //缓存，减少值交换
  //缓存不能够空间换时间，在1/4到3/4范围内插入还是会移动1/2的元素

  private val maxIdx: Int
    get() = if (minIdx == 0) list.size - 1 else minIdx - 1

  fun convertIdx(idx: Int): Int {
    return if (minIdx > list.size - idx - 1) idx + minIdx - list.size else idx + minIdx
  }

  fun minIdxRight(step: Int): Int {
    return if (minIdx + step > list.size - 1) minIdx + step - list.size else minIdx + step
  }

  fun minIdxLeft(step: Int): Int {
    return if (minIdx - step >= 0) minIdx - step else minIdx - step + list.size
  }

  fun insert(v: T) {
    if (v <= list[minIdx]) return
    if (v <= list[minIdxRight(1)]) {
      list[minIdx] = v
      return
    }
    if (v >= list[maxIdx]) {
      list[minIdx] = v
      minIdx = minIdxRight(1)
      return
    }

    var left = 0
    var right = list.size - 1
    var middle: Int
    while (left <= right) {
      middle = (left + right) / 2
      if (v < list[convertIdx(middle)]) {
        right = middle - 1
      } else {
        left = middle + 1
      }
    }
    //肯定可以再优化
    if (left > list.size / 2) {
      list[minIdx] = list[maxIdx]
      for (i in list.size - 1 downTo left) {
        list[convertIdx(i)] = list[convertIdx(i - 1)]
      }
      list[convertIdx(left)] = v
      minIdx = minIdxRight(1)
      return
    }
    for (i in 0 until left) {
      list[convertIdx(i)] = list[convertIdx(i + 1)]
    }
    list[convertIdx(left - 1)] = v
  }
}
