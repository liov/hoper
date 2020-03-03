package xyz.hoper.user.client;

import io.grpc.ManagedChannel;
import lombok.extern.log4j.Log4j2;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import xyz.hoper.protobuf.user.UserServiceGrpc;
import xyz.hoper.protobuf.user.UserServiceOuterClass;

import javax.annotation.PostConstruct;

@Component
@Log4j2
public class UserClient {

    @Value("${grpc.client.host}")
    private String host;
    @Value("${grpc.client.port}")
    private Integer port;

    @Autowired
    private GrpcClientMananer grpcClientMananer;

    @PostConstruct
    public void init() {
        call();
    }

    public void call(){
        ManagedChannel channel = grpcClientMananer.getChannel(host,port);
        UserServiceGrpc.UserServiceBlockingStub stub = UserServiceGrpc.newBlockingStub(channel);
        UserServiceOuterClass.GetReq request = UserServiceOuterClass.GetReq.newBuilder().setId(1).build();
        UserServiceOuterClass.GetRep reply = stub.getUser(request);
        log.info("time: "+ reply);
    }
}
