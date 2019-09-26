package xyz.hoper.client;

import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.springframework.stereotype.Component;

import java.util.concurrent.TimeUnit;

@Component
public class GrpcClientMananer {

    public ManagedChannel getChannel(String host, int port){
        ManagedChannel channel = ManagedChannelBuilder.forAddress(host,port)
                .usePlaintext()
                .disableRetry()
                .idleTimeout(5, TimeUnit.SECONDS)
                .build();
        return channel;
    }
}
