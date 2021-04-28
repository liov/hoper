package xyz.hoper.vertx.annotation

import kotlin.annotation.Retention

/**
 * Router API类 标识注解
 */
@Target(AnnotationTarget.ANNOTATION_CLASS, AnnotationTarget.CLASS)
@Retention(AnnotationRetention.RUNTIME)
annotation class RouteHandler(
        val value: String = "",
        val isOpen: Boolean = false,
        /**
         * 注册顺序，数字越大越先注册
         */
        val order: Int = 0)