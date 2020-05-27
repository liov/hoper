package xyz.hoper.vertx.annotation

import kotlin.annotation.Retention

/**
 * 异步服务
 *
 */
@Target(AnnotationTarget.ANNOTATION_CLASS, AnnotationTarget.CLASS)
@Retention(AnnotationRetention.RUNTIME)
annotation class AsyncServiceHandler 