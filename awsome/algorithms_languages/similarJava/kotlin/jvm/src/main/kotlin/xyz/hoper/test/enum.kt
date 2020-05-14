package xyz.hoper.test

import java.io.File
import java.util.function.BinaryOperator
import java.util.function.IntBinaryOperator

enum class Color(val rgb: Int) {
  RED(0xFF0000),
  GREEN(0x00FF00),
  BLUE(0x0000FF)
}

enum class ProtocolState {
  WAITING {
    override fun signal() = TALKING
  },

  TALKING {
    override fun signal() = WAITING
  };

  abstract fun signal(): ProtocolState
}


enum class IntArithmetics : BinaryOperator<Int>, IntBinaryOperator {
  PLUS {
    override fun apply(t: Int, u: Int): Int = t + u
  },
  TIMES {
    override fun apply(t: Int, u: Int): Int = t * u
  };

  override fun applyAsInt(t: Int, u: Int) = apply(t, u)
}

enum class RGB { RED, GREEN, BLUE }
inline fun <reified T : Enum<T>> printAllValues() {
  println(enumValues<T>().joinToString { it.name })
}


fun main(args: Array<String>){
  printAllValues<RGB>()
  println(enumValueOf<RGB>("RED"))
}

