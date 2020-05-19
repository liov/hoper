package xyz.hoper.test.oop

import kotlin.properties.Delegates
import kotlin.properties.ReadOnlyProperty
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

  println(lazyValue)
  println(lazyValue)

  MyUI()
}

//在每个委托属性的实现的背后，Kotlin 编译器都会生成辅助属性并委托给它
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
//lazy
val lazyValue: String by lazy(LazyThreadSafetyMode.PUBLICATION){
  println("computed!")
  "Hello"
}

//可观察属性 Observable
class DUser {
  var name: String by Delegates.observable("<no name>") {
    prop, old, new ->
    println("$old -> $new")
  }
}
//映射实例自身作为委托来实现委托属性
class DMUser(val map: Map<String, Any?>) {
  val name: String by map
  val age: Int     by map
}

fun example(computeFoo: () -> Foo) {
  val memoizedFoo by lazy(computeFoo)
//memoizedFoo 变量只会在第一次访问时计算。 如果 someCondition 失败，那么该变量根本不会计算。
/*  if (someCondition && memoizedFoo.isValid()) {
    memoizedFoo.doSomething()
  }*/
}

/*
class C {
  var prop: Type by MyDelegate()
}

// 这段是由编译器生成的相应代码：
class C {
  private val prop$delegate = MyDelegate()
  var prop: Type
    get() = prop$delegate.getValue(this, this::prop)
  set(value: Type) = prop$delegate.setValue(this, this::prop, value)
}*/


class ResourceDelegate<T>(val id: T) : ReadOnlyProperty<MyUI, T> {
  override fun getValue(thisRef: MyUI, property: KProperty<*>): T {
    println("$thisRef -> ${property.name}")
    return id
  }
}

class ResourceLoader<T>(rid: ResourceID<T>) {
  val _rid:T = rid.value
  operator fun provideDelegate(thisRef: MyUI, prop: KProperty<*>): ReadOnlyProperty<MyUI, T> {
    checkProperty(thisRef, prop.name)
    // 创建委托
    return ResourceDelegate(_rid)
  }

  private fun checkProperty(thisRef: MyUI, name: String) {
    if (name.startsWith("image")) {
      throw Exception("Image not supported")
    }
  }
}

class MyUI {
  fun <T> bindResource(id: ResourceID<T>): ResourceLoader<T> {
    return ResourceLoader(id)
  }

  val image by bindResource(ResourceID.image_id)
  val text by bindResource(ResourceID.text_id)
}

class ResourceID<T>(var value: T) {
  companion object {
    val image_id = ResourceID("image_id")
    val text_id = ResourceID(123)
  }
}
