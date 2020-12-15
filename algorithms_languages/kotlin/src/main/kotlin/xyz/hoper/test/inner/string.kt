package xyz.hoper.test.inner

fun main(args: Array<String>) {
  val s = StringBuilder()
  s.append('H')
  s.append('e')
  s.append('l')
  s.append('l')
  s.append('o')
  println(s.toString() == "Hello")
  val arr = CharArray(2)
  arr[0] = 'H'
  arrayDispatch(arr)
}
//传的是地址
fun arrayDispatch(arr: CharArray) {
  arr[1] = 'e'
  println(arr)
}
