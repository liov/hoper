package xyz.hoper.test

/**
 * 型变
 * Java 类型系统中最棘手的部分之一是通配符类型（参见 Java Generics FAQ）。
 * 而 Kotlin 中没有。 相反，它有两个其他的东西：声明处型变（declaration-site variance）与类型投影（type projections）。
 * 首先，让我们思考为什么 Java 需要那些神秘的通配符。
 * 在 《Effective Java》第三版 解释了该问题——第 31 条：利用有限制通配符来提升 API 的灵活性。
 *  首先，Java 中的泛型是不型变的，这意味着 List<String> 并不是 List<Object> 的子类型。
 * 为什么这样？ 如果 List 不是不型变的，它就没比 Java 的数组好到哪去，因为如下代码会通过编译然后导致运行时异常：
 * 我终于理解型变的意义了，回想起写go时 interface{}和[]interface{}的恶心
 * 协变（covariant）逆变（contravariance）
 *
 */
class Box<T>(t: T) {
  var value = t
}

interface Source<out T> {
  fun nextT(): T
}

fun demo(strs: Source<String>) {
  val objects: Source<Any> = strs // 这个没问题，因为 T 是一个 out-参数
  // ……
}

/**
 * 一般原则是：当一个类 C 的类型参数 T 被声明为 out 时，它就只能出现在 C 的成员的输出-位置，但回报是 C<Base> 可以安全地作为 C<Derived>的超类。
 * 简而言之，他们说类 C 是在参数 T 上是协变的，或者说 T 是一个协变的类型参数。 你可以认为 C 是 T 的生产者，而不是 T 的消费者。
 * 另外除了 out，Kotlin 又补充了一个型变注释：in。它使得一个类型参数逆变：只可以被消费而不可以被生产。逆变类型的一个很好的例子是 Comparable：
 * 我们相信 in 和 out 两词是自解释的（因为它们已经在 C# 中成功使用很长时间了）， 因此上面提到的助记符不是真正需要的，并且可以将其改写为更高的目标：
 */

interface Comparable<in T> {
  operator fun compareTo(other: T): Int
}

fun demo(x: Comparable<Number>) {
  x.compareTo(1.0) // 1.0 拥有类型 Double，它是 Number 的子类型
  // 因此，我们可以将 x 赋给类型为 Comparable <Double> 的变量
  val y: Comparable<Double> = x // OK！
}


fun copy(from: Array<out Any>, to: Array<Any>) {
  assert(from.size == to.size)
  for (i in from.indices)
    to[i] = from[i]
}

fun fill(dest: Array<in String>, value: String){}

fun <T> singletonList(item: T): List<T> {
  // ……
  return listOf(item)
}

/**
 * 星投影
 * 有时你想说，你对类型参数一无所知，但仍然希望以安全的方式使用它。 这里的安全方式是定义泛型类型的这种投影，该泛型类型的每个具体实例化将是该投影的子类型。
 * Kotlin 为此提供了所谓的星投影语法：
 * - 对于 Foo <out T : TUpper>，其中 T 是一个具有上界 TUpper 的协变类型参数，Foo <*> 等价于 Foo <out TUpper>。 这意味着当 T 未知时，你可以安全地从 Foo <*> 读取 TUpper 的值。
 * - 对于 Foo <in T>，其中 T 是一个逆变类型参数，Foo <*> 等价于 Foo <in Nothing>。 这意味着当 T 未知时，没有什么可以以安全的方式写入 Foo <*>。
 * - 对于 Foo <T : TUpper>，其中 T 是一个具有上界 TUpper 的不型变类型参数，Foo<*> 对于读取值时等价于 Foo<out TUpper> 而对于写值时等价于 Foo<in Nothing>。
 * 如果泛型类型具有多个类型参数，则每个类型参数都可以单独投影。 例如，如果类型被声明为 interface Function <in T, out U>，我们可以想象以下星投影：
 * - Function<*, String> 表示 Function<in Nothing, String>
 * - Function<Int, *> 表示 Function<Int, out Any?>；
 * - Function<*, *> 表示 Function<in Nothing, out Any?>。
 * 注意：星投影非常像 Java 的原始类型，但是安全。
 */


fun <T> T.basicToString(): String {  // 扩展函数
  // ……
  return "Basic"
}

//类rust
fun <T : Comparable<T>> sort(list: List<T>) {   }

fun <T> copyWhenGreater(list: List<T>, threshold: T): List<String>
  where T : CharSequence,
        T : Comparable<T> {
  return list.filter { it > threshold }.map { it.toString() }
}
