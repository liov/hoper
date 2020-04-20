package xyz.hoper.utils.vertx

import io.vertx.ext.web.Router
import xyz.hoper.utils.vertx.VertxUtil.vertxInstance

/**
 * 全局router单例
 */
object RouterUtil {
    var router: Router = Router.router(vertxInstance)
        private set
}