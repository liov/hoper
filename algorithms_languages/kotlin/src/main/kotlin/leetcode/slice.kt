package leetcode

class Slice<T> {
  var len: Int = 0
  var cap: Int = len
  lateinit var array: Array<T>

  constructor(len: Int = 0) {
    this.len = len
  }

  constructor(len: Int, cap: Int) {
    this.len = len
    this.cap = cap
  }
}
