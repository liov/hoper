package xyz.hoper.vertx.vertx.factory


import io.vertx.core.Handler
import io.vertx.core.http.HttpHeaders
import io.vertx.core.http.HttpMethod
import io.vertx.ext.web.Router
import io.vertx.ext.web.RoutingContext
import io.vertx.ext.web.handler.BodyHandler
import io.vertx.ext.web.handler.CorsHandler
import org.reflections.Reflections
import org.slf4j.LoggerFactory
import xyz.hoper.util.ReflectionUtil
import xyz.hoper.vertx.annotation.RouteHandler
import xyz.hoper.vertx.annotation.RouteMapping
import xyz.hoper.vertx.annotation.RouteMethod
import xyz.hoper.vertx.util.RouterUtil
import java.lang.reflect.InvocationTargetException
import java.lang.reflect.Method
import java.util.*
import java.util.stream.Collectors
import java.util.stream.Stream

/**
 * 创建Router对象
 */
class RouterHandlerFactory {

    companion object {
        private val LOGGER = LoggerFactory.getLogger(RouterHandlerFactory::class.java)

        @Volatile
        private lateinit var reflections: Reflections
        private const val GATEWAY_PREFIX = "/"
    }

    @Volatile
    private var gatewayPrefix = GATEWAY_PREFIX

    //String... routerScanAddress
    constructor(vararg routerScanAddress: String?) {
        Objects.requireNonNull(routerScanAddress, "The router package address scan is empty.")
        reflections = ReflectionUtil.getReflections(*routerScanAddress)
    }

    constructor(routerScanAddress: String, gatewayPrefix: String) {
        Objects.requireNonNull(routerScanAddress, "The router package address scan is empty.")
        reflections = ReflectionUtil.getReflections(routerScanAddress)
        this.gatewayPrefix = gatewayPrefix
    }

    /**
     * 扫描路由router并注册处理器handler
     */
    fun createRouter(): Router {
        val router = RouterUtil.router
        router.route().handler { ctx: RoutingContext ->
            //设置header
            LOGGER.debug("The HTTP service request address information ===>path:{}, uri:{}, method:{}",
                    ctx.request().path(), ctx.request().absoluteURI(), ctx.request().method())
            ctx.response().headers().add(HttpHeaders.CONTENT_TYPE, "application/json; charset=utf-8")
            ctx.response().headers().add(HttpHeaders.ACCESS_CONTROL_ALLOW_ORIGIN, "*")
            ctx.response().headers().add(HttpHeaders.ACCESS_CONTROL_ALLOW_METHODS, "POST, GET, OPTIONS, PUT, DELETE, HEAD")
            ctx.response().headers().add(HttpHeaders.ACCESS_CONTROL_ALLOW_HEADERS,
                    "X-PINGOTHER, Origin,Content-Type, Accept, X-Requested-With, Dev, Authorization, Version, Token")
            ctx.response().headers().add(HttpHeaders.ACCESS_CONTROL_MAX_AGE, "1728000")
            ctx.next()
        }
        val method: HashSet<HttpMethod?> = object : HashSet<HttpMethod?>() {
            init {
                add(HttpMethod.GET)
                add(HttpMethod.POST)
                add(HttpMethod.OPTIONS)
                add(HttpMethod.PUT)
                add(HttpMethod.DELETE)
                add(HttpMethod.HEAD)
            }
        }
        /* 添加跨域的方法 **/router.route().handler(CorsHandler.create("*").allowedMethods(method))
        router.route().handler(BodyHandler.create())

        //扫描处理器handler并注册到路由router，路由地址相同则比对两个路由权重顺序
        val sortedHandlers: List<Class<*>> = reflections.getTypesAnnotatedWith(RouteHandler::class.java).stream().sorted { m1, m2 ->
            val mapping1: RouteHandler = m1.getAnnotation(RouteHandler::class.java)
            val mapping2: RouteHandler = m2.getAnnotation(RouteHandler::class.java)
            mapping2.order.compareTo(mapping1.order)
        }.collect(Collectors.toList())
        for (handler in sortedHandlers) {
            registerNewHandler(router, handler)
        }
        return router
    }

    /**
     * 映射路由router到处理器handler
     */
    private fun registerNewHandler(router: Router, handler: Class<*>) {
        //默认api前缀
        var root = gatewayPrefix
        if (!root.startsWith("/")) {
            root = "/$root"
        }
        if (!root.endsWith("/")) {
            root = "$root/"
        }

        // 扫描RouteHandler注解的类
        if (handler.isAnnotationPresent(RouteHandler::class.java)) {
            val routeHandler: RouteHandler = handler.getAnnotation(RouteHandler::class.java)
            root += routeHandler.value
        }
        val methodList = Stream.of(*handler.methods).filter { method: Method -> method.isAnnotationPresent(RouteMapping::class.java) }.sorted { m1: Method, m2: Method ->
            val mapping1: RouteMapping = m1.getAnnotation(RouteMapping::class.java)
            val mapping2: RouteMapping = m2.getAnnotation(RouteMapping::class.java)
            mapping2.order.compareTo(mapping1.order)
        }.collect(Collectors.toList())
        try {
            val instance = handler.newInstance()
            for (method in methodList) {
                if (method.isAnnotationPresent(RouteMapping::class.java)) {
                    val mapping: RouteMapping = method.getAnnotation(RouteMapping::class.java)
                    val routeMethod: RouteMethod = mapping.method
                    var routeUrl: String
                    if (mapping.value.startsWith("/:")) {
                        routeUrl = method.name + mapping.value // 注意
                    } else {
                        routeUrl = if (mapping.value.endsWith(method.name)) mapping.value else if (mapping.isCover) mapping.value else mapping.value.toString() + method.name
                        if (routeUrl.startsWith("/")) {
                            routeUrl = routeUrl.substring(1)
                        }
                    }
                    val url = if (root.endsWith("/")) root + routeUrl else "$root/$routeUrl"
                    val methodHandler = method.invoke(instance) as Handler<RoutingContext>
                    LOGGER.info("Register New Handler -> {}:{}", routeMethod, url)
                    when (routeMethod) {
                        RouteMethod.GET -> router.get(url).handler(methodHandler)
                        RouteMethod.POST -> router.post(url).handler(methodHandler)
                        RouteMethod.PUT -> router.put(url).handler(methodHandler)
                        RouteMethod.DELETE -> router.delete(url).handler(methodHandler)
                        RouteMethod.ROUTE -> router.route(url).handler(methodHandler) // get、post、delete、put
                        else -> {
                        }
                    }
                }
            }
        } catch (ex: InstantiationException) {
            LOGGER.error("Obtain Handler Fail，Error details：{}", ex.message)
        } catch (ex: IllegalAccessException) {
            LOGGER.error("Obtain Handler Fail，Error details：{}", ex.message)
        } catch (ex: IllegalArgumentException) {
            LOGGER.error("Obtain Handler Fail，Error details：{}", ex.message)
        } catch (ex: InvocationTargetException) {
            LOGGER.error("Obtain Handler Fail，Error details：{}", ex.message)
        }
    }

}