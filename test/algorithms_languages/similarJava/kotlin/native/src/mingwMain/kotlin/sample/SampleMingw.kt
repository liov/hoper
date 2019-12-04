package sample

fun main(args: Array<String>) {
    val t1 = kotlin.system.getTimeMillis()

    for(i in 1..10000){
        aaa(i.toFloat())
    }

    val t2 = kotlin.system.getTimeMillis()

    println("java time: " + (t2 - t1).toString() + "ms")

}



fun  aaa(i:Float) {
    var a = i + 1
    var b = 2.3f
    val s = "abcdefkkbghisdfdfdsfds"

    if(a > b){
        ++a
    }else{
        b += 1
    }

    if(a == b) b += 1

    var c = (a * b  + a / b - a*a)

    var d = s.substring(0, s.indexOf("kkb")) + c.toString()

    print(d)
}
