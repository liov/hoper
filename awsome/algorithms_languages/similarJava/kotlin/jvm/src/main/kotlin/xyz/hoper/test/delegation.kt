package xyz.hoper.test

import kotlin.reflect.KProperty

interface DBase {
  val message: String
  fun print()
  fun printMessage()
  fun printMessageLine()
}

class BaseImpl(val x: Int) : DBase {
  override val message = "BaseImpl: x = $x"
  override fun print() { println(message) }
  override fun printMessage() { println(x) }
  override fun printMessageLine() { println(x) }
}

//Derived 的超类型列表中的 by-子句表示 b 将会在 Derived 中内部存储， 并且编译器将生成转发给 b 的所有 Base 的方法。
class DDerived(b: DBase) : DBase by b {
  // 在 b 的 `print` 实现中不会访问到这个属性
  override val message = "Message of Derived"
  override fun printMessage() { println("abc") }
}

fun main() {
  val b = BaseImpl(10)
  DDerived(b).printMessage()
  DDerived(b).printMessageLine()
  val derived = DDerived(b)
  derived.print()
  println(derived.message)
}

class Delegate {
  operator fun getValue(thisRef: Any?, property: KProperty<*>): String {
    return "$thisRef, thank you for delegating '${property.name}' to me!"
  }

  operator fun setValue(thisRef: Any?, property: KProperty<*>, value: String) {
    println("$value has been assigned to '${property.name}' in $thisRef.")
  }
}

class Example {
  var p: String by Delegate()
}
