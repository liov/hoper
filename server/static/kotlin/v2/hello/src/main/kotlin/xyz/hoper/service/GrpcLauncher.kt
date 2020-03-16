package xyz.hoper.service

import io.grpc.BindableService
import io.grpc.Server
import io.grpc.ServerBuilder
import lombok.extern.log4j.Log4j2
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
import xyz.hoper.client.HelloWorldClient
import java.io.IOException

@Component("grpcLauncher")
@Log4j2
class GrpcLauncher {
    private var server: Server? = null

    @Value("\${grpc.server.port}")
    private val grpcServerPort: Int = 8080

    val log : Logger = LoggerFactory.getLogger(GrpcLauncher::class.java)
    /**
     * GRPC 服务启动方法
     * @param grpcServiceBeanMap
     */
    fun grpcStart(grpcServiceBeanMap: Map<String?, Any>) {
        try {
            val serverBuilder = ServerBuilder.forPort(grpcServerPort)
            for (bean in grpcServiceBeanMap.values) {
                serverBuilder.addService(bean as BindableService)
                log.info(bean.javaClass.simpleName + " is regist in Spring Boot")
            }
            server = serverBuilder.build().start()
            log.info("grpc server is started at $grpcServerPort")
            server!!.awaitTermination()
            Runtime.getRuntime().addShutdownHook(Thread(Runnable { grpcStop() }))
        } catch (e: IOException) {
            e.printStackTrace()
        } catch (e: InterruptedException) {
            e.printStackTrace()
        }
    }

    private fun grpcStop() {
        if (server != null) {
            server!!.shutdownNow()
        }
    }
}
