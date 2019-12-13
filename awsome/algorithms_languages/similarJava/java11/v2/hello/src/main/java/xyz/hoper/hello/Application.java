package xyz.hoper.hello;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ConfigurableApplicationContext;
import xyz.hoper.hello.annotation.GrpcService;
import xyz.hoper.hello.service.GrpcLauncher;

import java.util.Map;

@SpringBootApplication
public class Application {
    public static void main(String[] args) {
        ConfigurableApplicationContext configurableApplicationContext = SpringApplication.run(Application.class, args);
        Map<String, Object> grpcServiceBeanMap =  configurableApplicationContext.getBeansWithAnnotation(GrpcService.class);
        GrpcLauncher grpcLauncher = configurableApplicationContext.getBean("grpcLauncher",GrpcLauncher.class);
        grpcLauncher.grpcStart(grpcServiceBeanMap);
    }
}
