package xyz.hoper.vertx

import io.vertx.core.AbstractVerticle
import io.vertx.core.Future
import io.vertx.core.Promise

class MainVerticle : AbstractVerticle() {

    override fun start(startPromise: Promise<Void>) {
        vertx
                .createHttpServer()
                .requestHandler { req ->
                    req.response()
                            .putHeader("content-type", "text/plain")
                            .end("测试!")
                }
                .listen(8888) { http ->
                    if (http.succeeded()) {
                        startPromise.complete()
                        println("HTTP server started on port 8888")
                    } else {
                        startPromise.fail(http.cause());
                    }
                }
    }
}
