package xyz.hoper.client;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.beans.factory.annotation.Value;

import java.util.concurrent.TimeUnit;

public class GrpcClientMananer {
    @Value("${grpc.client.host}")
    private String host;
    @Value("${grpc.client.port}")
    private Integer port;

    public ManagedChannel getChannel(){
        ManagedChannel channel = ManagedChannelBuilder.forAddress(host,port)
                .disableRetry()
                .idleTimeout(2, TimeUnit.SECONDS)
                .build();
        return channel;
    }
}
