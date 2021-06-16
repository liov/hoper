package xyz.hoper.test.coroutine

import kotlinx.coroutines.*
import kotlinx.coroutines.flow.*
import kotlin.system.measureTimeMillis

fun flowFoo(): Flow<Int> = flow { // 流构建器
  println("Flow started")
  for (i in 1..3) {
    delay(100) // 假装我们在这里做了一些有用的事情
    log("Emitting $i")
    emit(i) // 发送下一个值
  }
}

fun main() = runBlocking<Unit> {
  // 启动并发的协程以验证主线程并未阻塞
  launch {
    for (k in 1..3) {
      println("I'm not blocked $k")
      delay(100)
    }
  }
  println("Calling foo...")
  val flow = flowFoo()
  println("Calling collect...")
  flow.collect { value -> println(value) }
  println("Calling collect again...")
  withTimeoutOrNull(250) {
    flow.collect { value -> println(value) }
  }
  println("Done")

  (1..5).asFlow()
    .filter {
      println("Filter $it")
      it % 2 == 0
    }
    .map {
      println("Map $it")
      "string $it"
    }.collect {
      println("Collect $it")
    }

  flow.flowOn(Dispatchers.Default).collect { value ->
    log("Collected $value")  // 运行在指定上下文中
  }

  val time = measureTimeMillis {
    flow.buffer() // 缓冲发射项，无需等待
      .collect { value ->
        delay(300) // 假装我们花费 300 毫秒来处理它
        println(value)
      }
  }
  println("Collected in $time ms")


  val nums = (1..3).asFlow().onEach { delay(300) } // 发射数字 1..3，间隔 300 毫秒
  val strs = flowOf("one", "two", "three").onEach { delay(400) } // 每 400 毫秒发射一次字符串
  val startTime1 = System.currentTimeMillis() // 记录开始的时间
  nums.combine(strs) { a, b -> "$a -> $b" } // 使用“combine”组合单个字符串
    .collect { value -> // 收集并打印
      println("$value at ${System.currentTimeMillis() - startTime1} ms from start")
    }
  //一个包含流的流
  (1..3).asFlow().map { requestFlow(it) }
  val startTime2 = System.currentTimeMillis() // 记录开始时间
  (1..3).asFlow().onEach { delay(100) } // 每 100 毫秒发射一个数字
    .flatMapMerge { requestFlow(it) }
    .collect { value -> // 收集并打印
      println("$value at ${System.currentTimeMillis() - startTime2} ms from start")
    }

  try {
    flow.collect { value ->
      println(value)
      check(value <= 1) { "Collected $value" }
    }
  } catch (e: Throwable) {
    println("Caught $e")
  }

  flow.onCompletion { println("Done") }
    .map { value ->
    check(value <= 1) { "Crashed on $value" }
    "string $value"
  }.catch { e -> emit("Caught $e") } // 发射一个异常
    .collect { value -> println(value) }
}


fun requestFlow(i: Int): Flow<String> = flow {
  emit("$i: First")
  delay(500) // 等待 500 毫秒
  emit("$i: Second")
}
