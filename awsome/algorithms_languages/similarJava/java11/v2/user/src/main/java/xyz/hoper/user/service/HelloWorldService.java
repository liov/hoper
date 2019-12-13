package xyz.hoper.user.service;

import io.grpc.stub.StreamObserver;
import lombok.extern.slf4j.Slf4j;
import org.lognet.springboot.grpc.GRpcService;
import xyz.hoper.protobuf.GreeterGrpc;
import xyz.hoper.protobuf.HelloReply;
import xyz.hoper.protobuf.HelloRequest;
import xyz.hoper.user.interceptor.LogInterceptor;

import java.util.Date;

@Slf4j
@GRpcService(interceptors = {LogInterceptor.class}, applyGlobalInterceptors = false)
public class HelloWorldService extends GreeterGrpc.GreeterImplBase {

    @Override
    public void sayHello(HelloRequest request, StreamObserver<HelloReply> responseObserver) {
        HelloReply.Builder helloReplyOrBuilder = HelloReply.newBuilder();
        helloReplyOrBuilder.setTime(new Date().getTime());
        responseObserver.onNext(helloReplyOrBuilder.build());
        responseObserver.onCompleted();
    }
}
