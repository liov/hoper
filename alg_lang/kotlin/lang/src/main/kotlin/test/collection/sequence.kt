package test.collection

fun main(args: Array<String>) {
  val oddNumbers1 = generateSequence(1) { it + 2 } // `it` 是上一个元素
  println(oddNumbers1.take(5).toList())

  val oddNumbersLessThan10 = generateSequence(1) { if (it < 10) it + 2 else null }
  println(oddNumbersLessThan10.toList())

  val oddNumbers2 = sequence {
    yield(1)
    yieldAll(listOf(3, 5))
    yieldAll(generateSequence(7) { it + 2 })
  }
  println(oddNumbers2.take(5).toList())

  val words = "The quick brown fox jumps over the lazy dog".split(" ")
  val lengthsList = words.filter { println("filter: $it"); it.length > 3 }
    .map { println("length: ${it.length}"); it.length }
    .take(4)

  println("Lengths of first 4 words longer than 3 chars:")
  println(lengthsList)

  // 将列表转换为序列
  val wordsSequence = words.asSequence()

  val lengthsSequence = wordsSequence.filter { println("filter: $it"); it.length > 3 }
    .map { println("length: ${it.length}"); it.length }
    .take(4)

  println("Lengths of first 4 words longer than 3 chars")
// 末端操作：以列表形式获取结果。
  println(lengthsSequence.toList())
}
