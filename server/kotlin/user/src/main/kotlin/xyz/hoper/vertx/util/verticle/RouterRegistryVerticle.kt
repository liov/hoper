package xyz.hoper.vertx.util.verticle

import io.vertx.core.AbstractVerticle
import io.vertx.core.AsyncResult
import io.vertx.core.Handler
import io.vertx.core.Promise
import io.vertx.core.http.HttpServer
import io.vertx.core.http.HttpServerOptions
import io.vertx.core.http.HttpServerRequest
import io.vertx.ext.web.Router
import org.slf4j.LoggerFactory

/**
 * 路由发布
 */
class RouterRegistryVerticle : AbstractVerticle {
    private var port = HTTP_PORT
    private var server: HttpServer? = null
    private var router: Router

    constructor(router: Router) {
        this.router = router
    }

    constructor(router: Router, port: Int) {
        this.router = router
        if (port > 0) {
            this.port = port
        }
    }

    @Throws(Exception::class)
    override fun start(startFuture: Promise<Void>) {
        LOGGER.debug("To start listening to port {} ......", port)
        super.start()
        val options = HttpServerOptions().setMaxWebSocketFrameSize(MAX_WEBSOCKET_FRAME_SIZE).setPort(port)
        server = vertx.createHttpServer(options)
        server?.requestHandler(Handler { event: HttpServerRequest -> router.handle(event) })
        server?.listen(Handler { result: AsyncResult<HttpServer?> ->
            if (result.succeeded()) {
                startFuture.complete()
            } else {
                startFuture.fail(result.cause())
            }
        })
    }

    @Throws(Exception::class)
    override fun stop(stopFuture: Promise<Void>) {
        super.stop()
        if (server == null) {
            stopFuture.complete()
            return
        }
        server!!.close { result: AsyncResult<Void?> ->
            if (result.failed()) {
                stopFuture.fail(result.cause())
            } else {
                stopFuture.complete()
            }
        }
    }

    companion object {
        private val LOGGER = LoggerFactory.getLogger(RouterRegistryVerticle::class.java)

        // 前端发送给后端数据的最大值，默认值65536 (64KB)
        private const val MAX_WEBSOCKET_FRAME_SIZE = 65536
        private const val HTTP_PORT = 8080
    }
}