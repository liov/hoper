package xyz.hoper.vertx.util.verticle


import io.vertx.core.*
import io.vertx.serviceproxy.ServiceBinder
import org.slf4j.LoggerFactory
import xyz.hoper.vertx.annotation.AsyncServiceHandler
import xyz.hoper.util.ReflectionUtil
import xyz.hoper.util.SpringContextUtil
import xyz.hoper.vertx.util.VertxUtil
import java.util.*
import java.util.function.Consumer

class AsyncRegistVerticle(packageAddress: String) : AbstractVerticle() {
    private var packageAddress: String? = packageAddress

    init {
        Objects.requireNonNull(packageAddress, "given scan package address is empty")
    }

    @Throws(Exception::class)
    override fun start(startFuture: Promise<Void>) {
        val handlers: Set<Class<*>> = ReflectionUtil.getReflections(packageAddress).getTypesAnnotatedWith(AsyncServiceHandler::class.java)
        val binder = ServiceBinder(VertxUtil.vertxInstance)
        if (handlers.isNotEmpty()) {
            val ftList: MutableList<Future<*>> = ArrayList()
            handlers.forEach(Consumer { asyncService ->
                val pt: Handler<AsyncResult<Void>> = Promise.promise()
                val ft = pt as Future<*>
                try {
                    val asInstance = SpringContextUtil.getBean(asyncService)
                    val getAddressMethod = asyncService.getMethod("getAddress")
                    val address = getAddressMethod.invoke(asInstance) as String
                    val getAsyncInterfaceClassMethod = asyncService.getMethod("getAsyncInterfaceClass")
                    val clazz = getAsyncInterfaceClassMethod.invoke(asInstance) as Class<Any>
                    binder.setAddress(address).register(clazz, asInstance).completionHandler(pt)
                } catch (e: Exception) {
                    e.printStackTrace()
                }
                ftList.add(ft)
            })
            CompositeFuture.all(ftList).onComplete { ar: AsyncResult<CompositeFuture?> ->
                if (ar.succeeded()) {
                    LOGGER.info("All async services registered")
                    startFuture.complete()
                } else {
                    LOGGER.error(ar.cause().message)
                    startFuture.fail(ar.cause())
                }
            }
        }
    }

    companion object {
        private val LOGGER = LoggerFactory.getLogger(AsyncRegistVerticle::class.java)
    }
}