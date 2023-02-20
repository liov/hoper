package test

fun fail(message: String): Nothing {
  throw IllegalArgumentException(message)
}

class Div(val name: String?)

fun main(args: Array<String>) {
  test()
  println("done")
}

fun test(){
  val person = Div(null)
  val s = person.name ?: fail("Name required")
  println(s)
}
