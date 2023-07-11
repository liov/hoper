package test.oop

@JvmInline
value class Password(val value: String)

interface Printable {
  fun prettyPrint(): String
}

@JvmInline
value class Name(val s: String): Printable {
  val length: Int
    get() = s.length

  fun greet() {
    println("Hello, $s")
  }
  override fun prettyPrint(): String = "Let's $s!"
}

interface I

@JvmInline
value class Foo(val i: Int) : I

fun asInline(f: Foo) {}
fun <T> asGeneric(x: T) {}
fun asInterface(i: I) {}
fun asNullable(i: Foo?) {}

fun <T> id(x: T): T = x


typealias NameTypeAlias = String
@JvmInline
value class NameInlineClass(val s: String)

fun acceptString(s: String) {}
fun acceptNameTypeAlias(n: NameTypeAlias) {}
fun acceptNameInlineClass(p: NameInlineClass) {}

fun main() {
  val name = Name("Kotlin")
  name.greet() // `greet` 方法会作为一个静态方法被调用
  println(name.length) // 属性的 get 方法会作为一个静态方法被调用
  println(name.prettyPrint())

  val f = Foo(42)

  asInline(f)    // 拆箱操作: 用作 Foo 本身
  asGeneric(f)   // 装箱操作: 用作泛型类型 T
  asInterface(f) // 装箱操作: 用作类型 I
  asNullable(f)  // 装箱操作: 用作不同于 Foo 的可空类型 Foo?

  // 在下面这里例子中，'f' 首先会被装箱（当它作为参数传递给 'id' 函数时）然后又被拆箱（当它从'id'函数中被返回时）
  // 最后， 'c' 中就包含了被拆箱后的内部表达(也就是 '42')， 和 'f' 一样
  val c = id(f)

  val nameAlias: NameTypeAlias = ""
  val nameInlineClass: NameInlineClass = NameInlineClass("")
  val string: String = ""

  acceptString(nameAlias) // 正确: 传递别名类型的实参替代函数中基础类型的形参
  //acceptString(nameInlineClass) // 错误: 不能传递内联类的实参替代函数中基础类型的形参

  // And vice versa:
  acceptNameTypeAlias(string) // 正确: 传递基础类型的实参替代函数中别名类型的形参
  //acceptNameInlineClass(string) // 错误: 不能传递基础类型的实参替代函数中内联类类型的形参
}

// 在 JVM 平台上被表示为'public final void compute(int x)'
fun compute(x: Int) { }

// 同理，在 JVM 平台上也被表示为'public final void compute(int x)'！
fun compute(x: UInt) { }
