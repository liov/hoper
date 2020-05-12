package xyz.hoper.test

class Person(val name: String) {
  var children: MutableList<Person> = mutableListOf<Person>();
  constructor(name: String,parent: Person):this(name){
    parent.children.add(this)
  }
}

fun main(args: Array<String>){
  var p = Person("test")
  print(p.name)
}
