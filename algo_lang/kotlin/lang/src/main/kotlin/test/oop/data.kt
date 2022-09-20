package test.oop

data class User(val name: String, val age: Int)

data class DPerson(val name: String) {
  var age: Int = 0
}

fun main(args: Array<String>) {
  val jack = User(name = "Jack", age = 1)
  val olderJack = jack.copy(age = 2)
  val (name, age) = jack
  println("$name, ${olderJack.age} years of age")
}
