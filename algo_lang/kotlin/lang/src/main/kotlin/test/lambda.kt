package test

import java.util.concurrent.locks.Lock
import javax.swing.tree.TreeNode

fun main() {
  val stringPlus: (String, String) -> String = String::plus
  val intPlus: Int.(Int) -> Int = Int::plus

  println(stringPlus.invoke("<-", "->"))
  println(stringPlus("Hello, ", "world!"))

  println(intPlus.invoke(1, 1))
  println(intPlus(1, 2))
  println(2.intPlus(3)) // 类扩展调用

  val sum = { x: Int, y: Int -> x + y }
  //传递末尾的 lambda 表达式
  val items = listOf(1, 2, 3, 4, 5)
  val product = items.fold(1) { acc, e -> acc * e }
  //单个参数的隐式名称
  items.filter { it > 0 }
}

fun <T, R> Collection<T>.fold(
  initial: R,
  combine: (acc: R, nextElement: T) -> R
): R {
  var accumulator: R = initial
  for (element: T in this) {
    accumulator = combine(accumulator, element)
  }
  return accumulator
}

inline fun <T> lock(lock: Lock, body: () -> T): T {  return body() }
inline fun foo(inlined: () -> Unit, noinline notInlined: () -> Unit) {  }

inline fun inlined(block: () -> Unit) {
  println("hi!")
}

fun ffoo() {
  inlined {
    return // OK：该 lambda 表达式是内联的
  }
}

inline fun f(crossinline body: () -> Unit) {
  val f = object: Runnable {
    override fun run() = body()
  }
  // ……
}

inline fun <reified T> TreeNode.findParentOfType(): T? {
  var p = parent
  while (p != null && p !is T) {
    p = p.parent
  }
  return p as T?
}
