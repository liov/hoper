package test.oop

import java.io.File
import javax.naming.Context

@ExperimentalUnsignedTypes
class Person(val name: String) {
  var age = 0u
  var children: MutableList<Person> = mutableListOf<Person>()
  constructor(name: String,parent: Person):this(name){
    parent.children.add(this)
  }

  constructor(name: String, age: UInt):this(name){
    this.age = age
  }
  val firstProperty = "First property: $name".also(::println)

  init {
    println("First initializer block that prints ${name}")
  }

  val secondProperty = "Second property: ${name.length}".also(::println)

  init {
    println("Second initializer block that prints ${name.length}")
  }

}

@ExperimentalUnsignedTypes
fun main(args: Array<String>){
  var p = Person("test", 5u)
  print(p.age)
}

class DontCreateMe private constructor () { /*……*/ }
//继承
open class BaseClass(p: Int)

class DerivedClass(p: Int) : BaseClass(p)

open class View{
  constructor(ctx: Context)
}

class MyView : View {
  constructor(ctx: Context) : super(ctx)
}

interface Shape {
  val vertexCount: Int
  open fun draw() { /*……*/ }
  fun fill() { /*……*/ }
}

class Circle(override val vertexCount: Int = 1) : Shape {
  override fun draw() { /*……*/ }
}

open class Rectangle(override val vertexCount: Int = 4) : Shape {
  init { println("Initializing Base") }
  override fun draw() { /*……*/ }
  val borderColor: String get() = "black"
}

class FilledRectangle : Rectangle() {
  init { println("Initializing Derived") }
  final override fun draw() {
    super.draw()
    println("Filling the rectangle")
  }

  val fillColor: String get() = super.borderColor
  //在一个内部类中访问外部类的超类，可以通过由外部类名限定的 super 关键字来实现：super@Outer：
  inner class Filler {
    fun fill() { /* …… */ }
    fun drawAndFill() {
      super@FilledRectangle.draw() // 调用 Rectangle 的 draw() 实现
      fill()
      println("Drawn a filled rectangle with color ${super@FilledRectangle.borderColor}") // 使用 Rectangle 所实现的 borderColor 的 get()
    }
  }
}

interface Polygon {
  fun draw() { /* …… */ } // 接口成员默认就是“open”的
}

class Square() : Rectangle(), Polygon {
  // 编译器要求覆盖 draw()：
  override fun draw() {
    super<Rectangle>.draw() // 调用 Rectangle.draw()
    super<Polygon>.draw() // 调用 Polygon.draw()
  }
}

abstract class Triangle : Polygon {
  abstract override fun draw()
}

object DataProviderManager {
  fun registerDataProvider(provider: MyClass) {
  }

  private var _allDataProviders: List<MyClass>? = null

  val allDataProviders: Collection<MyClass>
    get() {
      if (_allDataProviders == null) {
        _allDataProviders = listOf()
      }
      return _allDataProviders ?:throw AssertionError("Set to null by another thread")
    }
}

interface Factory<T> {
  fun create(): T
}

class MyClass {
  companion object : Factory<MyClass> {
    override fun create(): MyClass = MyClass()
  }
  lateinit var subject: TestSubject

  class TestSubject {

  }
}


open class Outer {
  private val a = 1
  protected open val b = 2
  internal val c = 3
  val d = 4  // 默认 public

  protected class Nested {
    public val e: Int = 5
  }
}

class Subclass : Outer() {
  // a 不可见
  // b、c、d 可见
  // Nested 和 e 可见

  override val b = 5   // “b”为 protected
}

class Unrelated(o: Outer) {
  // o.a、o.b 不可见
  // o.c 和 o.d 可见（相同模块）
  // Outer.Nested 不可见，Nested::e 也不可见
}


typealias FileTable<K> = MutableMap<K, MutableList<File>>
