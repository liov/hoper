package test

const val SUBSYSTEM_DEPRECATED: String = "This subsystem is deprecated"

@Deprecated(SUBSYSTEM_DEPRECATED) fun foo() {  }

private fun bar() { } // 在 example.kt 内可见

public var bar: Int = 5 // 该属性随处可见
  private set         // setter 只在 example.kt 内可见

internal val baz = 6    // 相同模块内可见
