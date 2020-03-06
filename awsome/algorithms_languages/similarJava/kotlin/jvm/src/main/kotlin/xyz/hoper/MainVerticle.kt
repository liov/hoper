package xyz.hoper

import io.vertx.core.AbstractVerticle
import io.vertx.core.Promise
//vert.x像node，一个运行时。
// 略像go里的fasthttp，性能很高，只不过用的人少，fasthttp是因为非官方的实现，vert.x貌似正在火起来
//IDEA里直接运行：
//Main class  : io.vertx.core.Launcher
//Program args: run xyz.hoper.MainVerticle
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
