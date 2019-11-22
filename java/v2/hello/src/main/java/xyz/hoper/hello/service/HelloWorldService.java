package xyz.hoper.hello.service;

import io.grpc.stub.StreamObserver;
import xyz.hoper.hello.annotation.GrpcService;
import xyz.hoper.protobuf.GreeterGrpc;
import xyz.hoper.protobuf.HelloReply;
import xyz.hoper.protobuf.HelloRequest;

import java.util.Date;

@GrpcService
public class HelloWorldService extends GreeterGrpc.GreeterImplBase {

    @Override
    public void sayHello(HelloRequest request, StreamObserver<HelloReply> responseObserver) {
        HelloReply.Builder helloReplyOrBuilder = HelloReply.newBuilder();
        helloReplyOrBuilder.setTime(new Date().getTime());
        responseObserver.onNext(helloReplyOrBuilder.build());
        responseObserver.onCompleted();
    }
}
