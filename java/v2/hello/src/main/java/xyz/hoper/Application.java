package xyz.hoper;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.ConfigurableApplicationContext;
import xyz.hoper.annotation.GrpcService;
import xyz.hoper.server.GrpcLauncher;

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
