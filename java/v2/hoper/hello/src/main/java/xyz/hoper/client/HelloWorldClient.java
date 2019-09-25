package xyz.hoper.client;

import io.grpc.ManagedChannel;
import org.springframework.beans.factory.annotation.Autowired;

public class HelloWorldClient {
    @Autowired
    private GrpcClientMananer grpcClientMananer;

    public void call(){
        ManagedChannel channel = grpcClientMananer.getChannel();
        HelloWorld.NameRequestOrBuilder nameRequestOrBuilder = HelloWorld.NameRequest.newBuilder();
        ((HelloWorld.NameRequest.Builder) nameRequestOrBuilder).setName("Geek");
        HelloWorldServiceGrpc.HelloWorldServiceBlockingStub stub = HelloWorldServiceGrpc.newBlockingStub(channel);
        HelloWorld.EchoResponse echoResponse = stub.welcome(((HelloWorld.NameRequest.Builder) nameRequestOrBuilder).build());
        System.out.println(echoResponse.getMergename());
    }
}
