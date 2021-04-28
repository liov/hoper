package xyz.hoper.web

import io.vertx.core.Vertx
import io.vertx.core.VertxOptions
import io.vertx.core.eventbus.EventBusOptions
import io.vertx.ext.web.Router
import org.springframework.beans.factory.annotation.Value
import org.springframework.boot.runApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.context.event.ApplicationReadyEvent
import org.springframework.context.annotation.ComponentScan
import org.springframework.context.event.EventListener
import xyz.hoper.vertx.util.VertxUtil
import xyz.hoper.vertx.util.factory.RouterHandlerFactory
import xyz.hoper.vertx.util.verticle.DeployVertxServer
import java.io.IOException

@ComponentScan("xyz.hoper")
@SpringBootApplication
open class UserApplication {
    /**
     * web api所在包路径
     */
    @Value("\${web-api-packages}")
    private val webApiPackages: String? = null

    /**
     * 异步服务所在包路径
     */
    @Value("\${async-service-impl-packages}")
    private val asyncServiceImplPackages: String? = null

    /**
     * http服务器端口号
     */
    @Value("\${http-server-port}")
    private val httpServerPort = 8080

    /**
     * 工作线程池大小（可根据实际情况调整）
     */
    @Value("\${worker-pool-size}")
    private val workerPoolSize = 20

    /**
     * 异步服务实例数量（建议和CPU核数相同）
     */
    @Value("\${async-service-instances}")
    private val asyncServiceInstances = 4

    @EventListener
    fun deployVerticle(event: ApplicationReadyEvent?) {
        val eventBusOptions = EventBusOptions()
        //便于调试 设定超时等时间较长 生产环境建议适当调整
        eventBusOptions.connectTimeout = 1000 * 60 * 30
        val vertx: Vertx = Vertx.vertx(
                VertxOptions().setWorkerPoolSize(workerPoolSize)
                        .setEventBusOptions(eventBusOptions)
                        .setBlockedThreadCheckInterval(999999999L)
                        .setMaxEventLoopExecuteTime(Long.MAX_VALUE)
        )

        VertxUtil.init(vertx)
        try {
            val router: Router = RouterHandlerFactory(webApiPackages).createRouter()
            DeployVertxServer.startDeploy(router, asyncServiceImplPackages, httpServerPort, asyncServiceInstances)
        } catch (e: IOException) {
            e.printStackTrace()
        }
    }
}


fun main(args: Array<String>) {
    runApplication<UserApplication>(*args)
}
