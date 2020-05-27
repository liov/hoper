package xyz.hoper.hello

import org.springframework.boot.SpringApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import xyz.hoper.grpc.annotation.GrpcService
import xyz.hoper.hello.service.GrpcLauncher

@SpringBootApplication
class Application {
    fun main(args: Array<String>) {
        val configurableApplicationContext = SpringApplication.run(Application::class.java, *args)
        val grpcServiceBeanMap = configurableApplicationContext.getBeansWithAnnotation(GrpcService::class.java)
        val grpcLauncher: GrpcLauncher = configurableApplicationContext.getBean("grpcLauncher", GrpcLauncher::class.java)
        grpcLauncher.grpcStart(grpcServiceBeanMap)
    }
}