package test.coroutine

import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*

fun main() = runBlocking {
  val channel = Channel<Int>()
  launch {
    // 这里可能是消耗大量 CPU 运算的异步逻辑，我们将仅仅做 5 次整数的平方并发送
    for (x in 1..5) channel.send(x * x)
    channel.close() // 我们结束发送
  }
  // 这里我们打印了 5 次被接收的整数：
  for (y in channel) println(y)
  println("Done!")

  val squares = produceSquares()
  squares.consumeEach { println(it) }
  println("Done!")

  val numbers = produceNumbers() // 从 1 开始生成整数
  val squares2 = square(numbers) // 整数求平方
  repeat(5) {
    println(squares2.receive()) // 输出前五个
  }
  println("Done!") // 至此已完成
  coroutineContext.cancelChildren() // 取消子协程

  var cur = numbersFrom(2)
  repeat(10) {
    val prime = cur.receive()
    println(prime)
    cur = filter(cur, prime)
  }
  coroutineContext.cancelChildren() // 取消所有的子协程来让主协程结束

  val producer = produceNumbers()
  repeat(5) { launchProcessor(it, producer) }
  delay(950)
  producer.cancel() // 取消协程生产者从而将它们全部杀死

  val channel2 = Channel<String>()
  launch { sendString(channel2, "foo", 200L) }
  launch { sendString(channel2, "BAR!", 500L) }
  repeat(6) { // 接收前六个
    println(channel2.receive())
  }
  coroutineContext.cancelChildren() // 取消所有子协程来让主协程结束

  val table = Channel<Ball>() // 一个共享的 table（桌子）
  launch { player("ping", table) }
  launch { player("pong", table) }
  table.send(Ball(0)) // 乒乓球
  delay(1000) // 延迟 1 秒钟
  coroutineContext.cancelChildren() // 游戏结束，取消它们

  val tickerChannel = ticker(delayMillis = 100, initialDelayMillis = 0) //创建计时器通道
  var nextElement = withTimeoutOrNull(1) { tickerChannel.receive() }
  println("Initial element is available immediately: $nextElement") // no initial delay

  nextElement = withTimeoutOrNull(50) { tickerChannel.receive() } // all subsequent elements have 100ms delay
  println("Next element is not ready in 50 ms: $nextElement")

  nextElement = withTimeoutOrNull(50) { tickerChannel.receive() }
  println("Next element is ready in 100 ms: $nextElement")

  // 模拟大量消费延迟
  println("Consumer pauses for 150ms")
  delay(150)
  // 下一个元素立即可用
  nextElement = withTimeoutOrNull(1) { tickerChannel.receive() }
  println("Next element is available immediately after large consumer delay: $nextElement")
  // 请注意，`receive` 调用之间的暂停被考虑在内，下一个元素的到达速度更快
  nextElement = withTimeoutOrNull(50) { tickerChannel.receive() }
  println("Next element is ready in 50ms after consumer pause in 150ms: $nextElement")

  tickerChannel.cancel() // 表明不再需要更多的元素
}

fun CoroutineScope.produceSquares(): ReceiveChannel<Int> = produce {
  for (x in 1..5) send(x * x)
}

fun CoroutineScope.produceNumbers() = produce<Int> {
  var x = 1
  while (true) {
    send(x++) // 在流中开始从 1 生产无穷多个整数
    delay(100) // 等待 0.1 秒
  }

}

fun CoroutineScope.square(numbers: ReceiveChannel<Int>): ReceiveChannel<Int> = produce {
  for (x in numbers) send(x * x)
}

fun CoroutineScope.numbersFrom(start: Int) = produce<Int> {
  var x = start
  while (true) send(x++) // 开启了一个无限的整数流
}

fun CoroutineScope.filter(numbers: ReceiveChannel<Int>, prime: Int) = produce<Int> {
  for (x in numbers) if (x % prime != 0) send(x)
}

fun CoroutineScope.launchProcessor(id: Int, channel: ReceiveChannel<Int>) = launch {
  for (msg in channel) {
    println("Processor #$id received $msg")
  }
}

suspend fun sendString(channel: SendChannel<String>, s: String, time: Long) {
  while (true) {
    delay(time)
    channel.send(s)
  }
}

data class Ball(var hits: Int)

suspend fun player(name: String, table: Channel<Ball>) {
  for (ball in table) { // 在循环中接收球
    ball.hits++
    println("$name $ball")
    delay(300) // 等待一段时间
    table.send(ball) // 将球发送回去
  }
}
