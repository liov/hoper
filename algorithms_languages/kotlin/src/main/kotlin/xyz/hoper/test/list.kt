package xyz.hoper.test

fun main() {
  val list1 = listOf(1, 2, 3)
  val list2 = listOf("1", "2", "3")
  val list3 = mutableListOf<Any>()
  list3.addAll(list1)
  list3.addAll(list2)
  println(list3)
  for(i in list3.indices){
    if (list3[i] is Int)
      print(list3[i] as Int)
    else if (list3[i] is String){
      print(list3[i] as String)
    }
  }
}
