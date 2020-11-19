package xyz.hoper.test.coroutine

import kotlinx.coroutines.*
import kotlin.system.measureTimeMillis

suspend fun doSomethingUsefulOne(): Int {
  delay(1000L) // 假设我们在这里做了一些有用的事
  return 13
}

suspend fun doSomethingUsefulTwo(): Int {
  delay(1000L) // 假设我们在这里也做了一些有用的事
  return 29
}

fun main() = runBlocking<Unit> {
  val time1 = measureTimeMillis {
    val one = async { doSomethingUsefulOne() }
    val two = async { doSomethingUsefulTwo() }
    println("The answer is ${one.await() + two.await()}")
  }
  println("Completed in $time1 ms")

  val time2 = measureTimeMillis {
    val one = async(start = CoroutineStart.LAZY) { doSomethingUsefulOne() }
    val two = async(start = CoroutineStart.LAZY) { doSomethingUsefulTwo() }
    // 执行一些计算
    one.start() // 启动第一个
    two.start() // 启动第二个
    println("The answer is ${one.await() + two.await()}")
  }
  println("Completed in $time2 ms")

  val time3 = measureTimeMillis {
    // 我们可以在协程外面启动异步执行
    val one = somethingUsefulOneAsync()
    val two = somethingUsefulTwoAsync()
    // 但是等待结果必须调用其它的挂起或者阻塞
    // 当我们等待结果的时候，这里我们使用 `runBlocking { …… }` 来阻塞主线程
    runBlocking {
      println("The answer is ${one.await() + two.await()}")
    }
  }
  println("Completed in $time3 ms")

  val time4 = measureTimeMillis {
    println("The answer is ${concurrentSum()}")
  }
  println("Completed in $time4 ms")
  //异常取消
  try {
    failedConcurrentSum()
  } catch(e: ArithmeticException) {
    println("Computation failed with ArithmeticException")
  }
  //
  launch { // 运行在父协程的上下文中，即 runBlocking 主协程
    println("main runBlocking      : I'm working in thread ${Thread.currentThread().name}")
  }
  launch(Dispatchers.Unconfined) { // 不受限的——将工作在主线程中
    println("Unconfined            : I'm working in thread ${Thread.currentThread().name}")
  }
  launch(Dispatchers.Default) { // 将会获取默认调度器
    println("Default               : I'm working in thread ${Thread.currentThread().name}")
  }
  launch(newSingleThreadContext("MyOwnThread")) { // 将使它获得一个新的线程
    println("newSingleThreadContext: I'm working in thread ${Thread.currentThread().name}")
  }

  launch(Dispatchers.Unconfined) { // 非受限的——将和主线程一起工作
    println("Unconfined      : I'm working in thread ${Thread.currentThread().name}")
    delay(500)
    println("Unconfined      : After delay in thread ${Thread.currentThread().name}")
  }
  launch { // 父协程的上下文，主 runBlocking 协程
    println("main runBlocking: I'm working in thread ${Thread.currentThread().name}")
    delay(1000)
    println("main runBlocking: After delay in thread ${Thread.currentThread().name}")
  }

  newSingleThreadContext("Ctx1").use { ctx1 ->
    newSingleThreadContext("Ctx2").use { ctx2 ->
      runBlocking(ctx1) {
        log("Started in ctx1")
        withContext(ctx2) {
          log("Working in ctx2")
        }
        log("Back to ctx1")
      }
    }
  }

  // 启动一个协程来处理某种传入请求（request）
  val request = launch {
    // 孵化了两个子作业, 其中一个通过 GlobalScope 启动
    GlobalScope.launch {
      println("job1: I run in GlobalScope and execute independently!")
      delay(1000)
      println("job1: I am not affected by cancellation of the request")
    }
    // 另一个则承袭了父协程的上下文
    launch {
      delay(100)
      println("job2: I am a child of the request coroutine")
      delay(1000)
      println("job2: I will not execute this line if my parent request is cancelled")
    }
  }
  delay(500)
  request.cancel() // 取消请求（request）的执行
  delay(1000) // 延迟一秒钟来看看发生了什么
  println("main: Who has survived request cancellation?")

  //-Dkotlinx.coroutines.debug
  log("Started main coroutine")
// 运行两个后台值计算
  val v1 = async(CoroutineName("v1coroutine")) {
    delay(500)
    log("Computing v1")
    252
  }
  val v2 = async(CoroutineName("v2coroutine")) {
    delay(1000)
    log("Computing v2")
    6
  }
  log("The answer for v1 / v2 = ${v1.await() / v2.await()}")

  val activity = Activity()
  activity.doSomething() // 运行测试函数
  println("Launched coroutines")
  delay(500L) // 延迟半秒钟
  println("Destroying activity!")
  activity.destroy() // 取消所有的协程
  delay(1000) // 为了在视觉上确认它们没有工作

  val threadLocal = ThreadLocal<String?>() // 声明线程局部变量
  threadLocal.set("main")
  println("Pre-main, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
  val job = launch(Dispatchers.Default + threadLocal.asContextElement(value = "launch")) {
    println("Launch start, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
    yield()
    println("After yield, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
  }
  job.join()
  println("Post-main, current thread: ${Thread.currentThread()}, thread local value: '${threadLocal.get()}'")
}

fun log(msg: String) = println("[${Thread.currentThread().name}] $msg")

// somethingUsefulOneAsync 函数的返回值类型是 Deferred<Int>
fun somethingUsefulOneAsync() = GlobalScope.async {
  doSomethingUsefulOne()
}

// somethingUsefulTwoAsync 函数的返回值类型是 Deferred<Int>
fun somethingUsefulTwoAsync() = GlobalScope.async {
  doSomethingUsefulTwo()
}

suspend fun concurrentSum(): Int = coroutineScope {
  val one = async { doSomethingUsefulOne() }
  val two = async { doSomethingUsefulTwo() }
  one.await() + two.await()
}
//如果其中一个子协程（即 two）失败，第一个 async 以及等待中的父协程都会被取消：
suspend fun failedConcurrentSum(): Int = coroutineScope {
  val one = async<Int> {
    try {
      delay(Long.MAX_VALUE) // 模拟一个长时间的运算
      42
    } finally {
      println("First child was cancelled")
    }
  }
  val two = async<Int> {
    println("Second child throws an exception")
    throw ArithmeticException()
  }
  one.await() + two.await()
}

class Activity {
  private val mainScope = CoroutineScope(Dispatchers.Default)

  fun destroy() {
    mainScope.cancel()
  }

  // 在 Activity 类中
  fun doSomething() {
    // 在示例中启动了 10 个协程，且每个都工作了不同的时长
    repeat(10) { i ->
      mainScope.launch {
        delay((i + 1) * 200L) // 延迟 200 毫秒、400 毫秒、600 毫秒等等不同的时间
        println("Coroutine $i is done")
      }
    }
  }
} // Activity 类结束
