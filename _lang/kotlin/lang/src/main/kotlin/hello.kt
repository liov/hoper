fun main() {
    for(i in 1..10000){
        aaa(i.toFloat())
    }
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

    val c = (a * b  + a / b - a*a)

    val d = s.substring(0, s.indexOf("kkb")) + c.toString()

    println(d)
}
