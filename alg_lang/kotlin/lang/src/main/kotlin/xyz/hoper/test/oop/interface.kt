package xyz.hoper.test.oop

interface MyInterface {
  val prop: Int // 抽象的

  val propertyWithImplementation: String
    get() = "foo"

  fun foo() {
    print(prop)
  }
}

class Child : MyInterface {
  override val prop: Int = 29
}

interface Named {
  val name: String
}

interface People : Named {
  val firstName: String
  val lastName: String

  override val name: String get() = "$firstName $lastName"
}

data class Employee(
  // 不必实现“name”
  override val firstName: String,
  override val lastName: String,
  val position: String
) : People

interface A {
  fun foo() { print("A") }
  fun bar()
}

interface B {
  fun foo() { print("B") }
  fun bar() { print("bar") }
}

class C : A {
  override fun bar() { print("bar") }
}

class D : A, B {
  override fun foo() {
    super<A>.foo()
    super<B>.foo()
  }

  override fun bar() {
    super<B>.bar()
  }
}
