package xyz.hoper.vertx.util

import io.vertx.ext.web.Router

/**
 * 全局router单例
 */
object RouterUtil {
    var router: Router = Router.router(VertxUtil.vertxInstance)
        private set
}