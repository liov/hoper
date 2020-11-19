package xyz.hoper.test.coroutine

import kotlinx.coroutines.*


fun main() = runBlocking { // this: CoroutineScope
  launch {
    delay(200L)
    println("Task from runBlocking")
  }

  coroutineScope { // 创建一个协程作用域
    launch {
      delay(500L)
      println("Task from nested launch")
    }

    delay(100L)
    println("Task from coroutine scope") // 这一行会在内嵌 launch 之前输出
  }

  println("Coroutine scope is over") // 这一行在内嵌 launch 执行完毕后才输出
  //取消协成
  val job = launch {
    repeat(1000) { i ->
      println("job: I'm sleeping $i ...")
      delay(500L)
    }
  }
  delay(1300L) // 延迟一段时间
  println("main: I'm tired of waiting!")
  job.cancel() // 取消该作业
  job.join() // 等待作业执行结束
  println("main: Now I can quit.")
}

suspend fun doWorld() {
  delay(1000L)
  println("World!")
}
