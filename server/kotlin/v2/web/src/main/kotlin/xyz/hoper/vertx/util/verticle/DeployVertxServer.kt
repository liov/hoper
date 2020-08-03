package xyz.hoper.vertx.util.verticle

import io.vertx.core.DeploymentOptions
import io.vertx.ext.web.Router
import org.slf4j.LoggerFactory
import xyz.hoper.vertx.util.VertxUtil
import java.io.IOException

/**
 * 注册vertx相关服务,并发布
 */
object DeployVertxServer {
    private val LOGGER = LoggerFactory.getLogger(DeployVertxServer::class.java)

    /**
     *
     * @param router
     * @param asyncServiceImplPackages
     * @param port
     * @param asyncServiceInstances 确定启动多少个程式实例
     * @throws IOException
     */
    @Throws(IOException::class)
    fun startDeploy(router: Router?, asyncServiceImplPackages: String?, port: Int, asyncServiceInstances: Int) {
        var asyncServiceInstances = asyncServiceInstances
        LOGGER.debug("Start Deploy....")
        LOGGER.debug("Start registry router....")
        VertxUtil.vertxInstance?.deployVerticle(router?.let { RouterRegistryVerticle(it, port) })
        LOGGER.debug("Start registry service....")
        if (asyncServiceInstances < 1) {
            asyncServiceInstances = 1
        }
        for (i in 0 until asyncServiceInstances) {
            VertxUtil.vertxInstance?.deployVerticle(AsyncRegistVerticle(asyncServiceImplPackages!!), DeploymentOptions().setWorker(true))
        }
    }
}