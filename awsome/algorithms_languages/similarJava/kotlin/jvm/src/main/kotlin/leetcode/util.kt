package leetcode

/***
 * 打印排序后数组
 */
fun printArrSorted(arr: IntArray) {
  arr.sort()
  printArr(arr)
}

/***
 * 打印数组
 */
fun printArr(arr: IntArray) {
  println(arr.joinToString(prefix="[",postfix="]"))
}

/***
 * 值交换
 */
//抽象能力太弱，明明能为所有实现getset操作符的类实现swap
fun IntArray.swap(i: Int, j: Int) {
  val tmp = this[i]
  this[i] = this[j]
  this[j] = tmp
}
