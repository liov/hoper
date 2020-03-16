package xyz.hoper

import org.springframework.boot.SpringApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import xyz.hoper.annotation.GrpcService
import xyz.hoper.service.GrpcLauncher

@SpringBootApplication
class Application {
    fun main(args: Array<String>) {
        val configurableApplicationContext = SpringApplication.run(Application::class.java, *args)
        val grpcServiceBeanMap = configurableApplicationContext.getBeansWithAnnotation(GrpcService::class.java)
        val grpcLauncher: GrpcLauncher = configurableApplicationContext.getBean("grpcLauncher", GrpcLauncher::class.java)
        grpcLauncher.grpcStart(grpcServiceBeanMap)
    }
}