package xyz.hoper.test.collection

//集合的基本操作，rust中是对实现迭代器类型的操作map,filter,reduce,take,zip,
fun main(args: Array<String>) {
  val numbers1 = listOf("one", "two", "three", "four")
  val filterResults = mutableListOf<String>()  // 目标对象
  numbers1.filterTo(filterResults) { it.length > 3 }
  numbers1.filterIndexedTo(filterResults) { index, _ -> index == 0 }
  println(filterResults) // 包含两个操作的结果

  // 将数字直接过滤到新的哈希集中，
// 从而消除结果中的重复项
  val result = numbers1.mapTo(HashSet()) { it.length }
  println("distinct item lengths are $result")

  val numbers2 = setOf(1, 2, 3)
  println(numbers2.mapNotNull { if ( it == 2) null else it * 3 })
  println(numbers2.mapIndexedNotNull { idx, value -> if (idx == 0) null else value * idx })

  val colors = listOf("red", "brown", "grey")
  val animals = listOf("fox", "bear", "wolf")

  println(colors.zip(animals) { color, animal -> "The ${animal.capitalize()} is $color"})

  val numberSets = listOf(setOf(1, 2, 3), setOf(4, 5, 6), setOf(1, 2))
  println(numberSets.flatten())

  val numbers = listOf("one", "two", "three", "four")

  println(numbers)
  println(numbers.joinToString())

  val listString = StringBuffer("The list of numbers: ")
  numbers.joinTo(listString)
  println(listString)
  println(numbers.joinToString(separator = " | ", prefix = "start: ", postfix = ": end"))
}
