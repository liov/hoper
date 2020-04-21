package xyz.hoper.vertx.util

import io.vertx.core.Vertx
import java.util.*

/**
 * 全局vertx单例
 */
object VertxUtil {
    private var singletonVertx: Vertx? = null
    fun init(vertx: Vertx?) {
        Objects.requireNonNull(vertx, "未初始化Vertx")
        singletonVertx = vertx
    }

    val vertxInstance: Vertx?
        get() {
            Objects.requireNonNull(singletonVertx, "未初始化Vertx")
            return singletonVertx
        }
}