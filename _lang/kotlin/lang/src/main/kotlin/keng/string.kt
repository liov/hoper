package keng

fun main(){
  val v = "test"
  println("\${v}")//支持反斜杠
  println("""\${v}""")//原始字符串不支持反斜杠
}
