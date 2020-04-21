package xyz.hoper.vertx.util

import io.vertx.core.Vertx
import io.vertx.serviceproxy.ServiceProxyBuilder

object AsyncServiceUtil {
    fun <T> getAsyncServiceInstance(asClazz: Class<T>, vertx: Vertx?): T {
        return ServiceProxyBuilder(vertx).setAddress(asClazz.name).build(asClazz)
    }

    fun <T> getAsyncServiceInstance(asClazz: Class<T>): T {
        return ServiceProxyBuilder(VertxUtil.vertxInstance).setAddress(asClazz.name).build(asClazz)
    }
}