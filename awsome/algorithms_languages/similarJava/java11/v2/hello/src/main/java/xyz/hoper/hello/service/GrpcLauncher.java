package xyz.hoper.hello.service;

import io.grpc.BindableService;
import io.grpc.Server;
import io.grpc.ServerBuilder;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.io.IOException;
import java.util.Map;

@Component("grpcLauncher")
@Log4j2
public class GrpcLauncher {

    private Server server;

    @Value("${grpc.server.port}")
    private Integer grpcServerPort;
    /**
     * GRPC 服务启动方法
     * @param grpcServiceBeanMap
     */
    public void grpcStart(Map<String, Object> grpcServiceBeanMap) {
        try{
            ServerBuilder serverBuilder = ServerBuilder.forPort(grpcServerPort);
            for (Object bean : grpcServiceBeanMap.values()){
                serverBuilder.addService((BindableService) bean);
                log.info(bean.getClass().getSimpleName() + " is regist in Spring Boot");
            }
            server = serverBuilder.build().start();
            log.info("grpc server is started at " +  grpcServerPort);
            server.awaitTermination();
            Runtime.getRuntime().addShutdownHook(new Thread(()-> grpcStop()));
        } catch (IOException | InterruptedException e){
            e.printStackTrace();
        }
    }
    private void grpcStop(){
        if (server != null){
            server.shutdownNow();
        }
    }
}
