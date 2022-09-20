package test.coroutine

import kotlinx.coroutines.*
import kotlinx.coroutines.channels.*
import kotlin.system.*
import kotlinx.coroutines.sync.*
import java.util.concurrent.atomic.*

suspend fun massiveRun(action: suspend () -> Unit) {
  val n = 100  // 启动的协程数量
  val k = 1000 // 每个协程重复执行同一动作的次数
  val time = measureTimeMillis {
    coroutineScope { // 协程的作用域
      repeat(n) {
        launch {
          repeat(k) { action() }
        }
      }
    }
  }
  println("Completed ${n * k} actions in $time ms")
}

var counter = AtomicInteger()
val counterContext = newSingleThreadContext("CounterContext")
var count = 0
val mutex = Mutex()

fun main() = runBlocking<Unit> {
  withContext(Dispatchers.Default) {
    massiveRun {
      counter.incrementAndGet()
    }
  }
  println("Counter = $counter")

  withContext(Dispatchers.Default) {
    massiveRun {
      // 将每次自增限制在单线程上下文中
      withContext(counterContext) {
        count++
      }
    }
  }
  println("Counter = $count")

  // 将一切都限制在单线程上下文中
  withContext(counterContext) {
    massiveRun {
      count++
    }
  }
  println("Counter = $count")

  withContext(Dispatchers.Default) {
    massiveRun {
      // 用锁保护每次自增
      mutex.withLock {
        count++
      }
    }
  }
  println("Counter = $count")

  val counter = counterActor() // 创建该 actor
  withContext(Dispatchers.Default) {
    massiveRun {
      counter.send(IncCounter)
    }
  }
  // 发送一条消息以用来从一个 actor 中获取计数值
  val response = CompletableDeferred<Int>()
  counter.send(GetCounter(response))
  println("Counter = ${response.await()}")
  counter.close() // 关闭该actor
}

// 计数器 Actor 的各种类型
sealed class CounterMsg
object IncCounter : CounterMsg() // 递增计数器的单向消息
class GetCounter(val response: CompletableDeferred<Int>) : CounterMsg() // 携带回复的请求

// 这个函数启动一个新的计数器 actor
fun CoroutineScope.counterActor() = actor<CounterMsg> {
  var counter = 0 // actor 状态
  for (msg in channel) { // 即将到来消息的迭代器
    when (msg) {
      is IncCounter -> counter++
      is GetCounter -> msg.response.complete(counter)
    }
  }
}
