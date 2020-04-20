package xyz.hoper.utils.vertx

import io.vertx.core.Vertx
import io.vertx.serviceproxy.ServiceProxyBuilder
import xyz.hoper.utils.vertx.VertxUtil.vertxInstance

object AsyncServiceUtil {
    fun <T> getAsyncServiceInstance(asClazz: Class<T>, vertx: Vertx?): T {
        return ServiceProxyBuilder(vertx).setAddress(asClazz.name).build(asClazz)
    }

    fun <T> getAsyncServiceInstance(asClazz: Class<T>): T {
        return ServiceProxyBuilder(vertxInstance).setAddress(asClazz.name).build(asClazz)
    }
}