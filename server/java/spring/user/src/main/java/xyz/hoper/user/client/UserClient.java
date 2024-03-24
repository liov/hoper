package xyz.hoper.user.client;


import com.google.protobuf.Empty;
import io.grpc.StatusRuntimeException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.annotation.PostConstruct;

import xyz.hoper.protobuf.user.UserServiceGrpc;


@Component
class UserClient {
    @Value("${grpc.client.host}")
    private String host;
    @Value("${grpc.client.port}")
    private Integer port = 8080;

    @Autowired
    private GrpcClientManager grpcClientManager;

    Logger log = LoggerFactory.getLogger(UserClient.class);

    @PostConstruct
    void init() {
        call();
    }

    void call() {
        var channel = grpcClientManager.getChannel(host, port);
        var stub = UserServiceGrpc.newBlockingStub(channel);
        var request = Empty.newBuilder().build();
        try {
            var reply = stub.verifyCode(request);
            log.info("time: $reply");
        } catch (StatusRuntimeException e) {
            log.error(e.getMessage());
        }
    }

}