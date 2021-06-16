package keng

fun list(){
  val list = ArrayList<Int>()
  list.add(1)
  list.remove(list.lastIndex)
  list.removeAt(list.lastIndex)
  //自动重载参数object
  list.remove(1)
}
