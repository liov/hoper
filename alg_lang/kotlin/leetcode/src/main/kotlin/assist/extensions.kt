package assist

//扩展
fun <T> MutableList<T>.swap(index1: Int, index2: Int) {
  val tmp = this[index1] // “this”对应该列表
  this[index1] = this[index2]
  this[index2] = tmp
}

fun <T> Array<T>.swap(index1: Int, index2: Int) {
  val tmp = this[index1] // “this”对应该列表
  this[index1] = this[index2]
  this[index2] = tmp
}

val <T> List<T>.lastIndex: Int
  get() = size - 1

class Host(val hostname: String) {
  fun printHostname() { print(hostname) }
}

class Connection(val host: Host, val port: Int) {
  fun printPort() { print(port) }

  fun Host.printConnectionString() {
    printHostname()   // 调用 Host.printHostname()
    print(":")
    printPort()   // 调用 Connection.printPort()
    toString()         // 调用 Host.toString()
    this@Connection.toString()  // 调用 Connection.toString()
  }

  fun connect() {
    /*……*/
    host.printConnectionString()   // 调用扩展函数
  }
}

open class Base { }

class Derived : Base() { }

open class BaseCaller {
  open fun Base.printFunctionInfo() {
    println("Base extension function in BaseCaller")
  }

  open fun Derived.printFunctionInfo() {
    println("Derived extension function in BaseCaller")
  }

  fun call(b: Base) {
    b.printFunctionInfo()   // 调用扩展函数
  }
}

class DerivedCaller: BaseCaller() {
  override fun Base.printFunctionInfo() {
    println("Base extension function in DerivedCaller")
  }

  override fun Derived.printFunctionInfo() {
    println("Derived extension function in DerivedCaller")
  }
}

fun main() {
  BaseCaller().call(Base())   // “Base extension function in BaseCaller”
  DerivedCaller().call(Base())  // “Base extension function in DerivedCaller”——分发接收者虚拟解析
  DerivedCaller().call(Derived())  // “Base extension function in DerivedCaller”——扩展接收者静态解析
}
