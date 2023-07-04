package design

import design.CustomerFactory.getCustomer

/**
在空对象模式（Null Object Pattern）中，一个空对象取代 NULL 对象实例的检查。Null 对象不是检查空值，而是反应一个不做任何动作的关系。这样的 Null 对象也可以在数据不可用的时候提供默认的行为。

在空对象模式中，我们创建一个指定各种要执行的操作的抽象类和扩展该类的实体类，还创建一个未对该类做任何实现的空对象类，该空对象类将无缝地使用在需要检查空值的地方。
  */
abstract class AbstractCustomer {
  open var name: String? = null
  abstract fun isNil(): Boolean
}
class RealCustomer(name: String?) : AbstractCustomer() {
  override fun isNil(): Boolean {
    return false
  }
}

class NullCustomer : AbstractCustomer() {
  override var name: String?
    get() = "Not Available in Customer Database"
    set(name) {
      super.name = name
    }

  override fun isNil(): Boolean {
    return true
  }
}

object CustomerFactory {
  val names = arrayOf("Rob", "Joe", "Julie")
  fun getCustomer(name: String?): AbstractCustomer {
    for (i in names.indices) {
      if (names[i].equals(name, ignoreCase = true)) {
        return RealCustomer(name)
      }
    }
    return NullCustomer()
  }
}


fun main() {
  val customer1 = getCustomer("Rob")
  val customer2 = getCustomer("Bob")
  val customer3 = getCustomer("Julie")
  val customer4 = getCustomer("Laura")
  println("Customers")
  println(customer1.name)
  println(customer2.name)
  println(customer3.name)
  println(customer4.name)
}
