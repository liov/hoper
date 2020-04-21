package xyz.hoper.hello.service

import io.grpc.stub.StreamObserver
import xyz.hoper.grpc.annotation.GrpcService
import xyz.hoper.protobuf.GreeterGrpc.GreeterImplBase
import xyz.hoper.protobuf.HelloReply
import xyz.hoper.protobuf.HelloRequest
import java.util.*

@GrpcService
class HelloWorldService : GreeterImplBase() {
    override fun sayHello(request: HelloRequest, responseObserver: StreamObserver<HelloReply>) {
        val helloReplyOrBuilder = HelloReply.newBuilder()
        helloReplyOrBuilder.time = Date().time
        responseObserver.onNext(helloReplyOrBuilder.build())
        responseObserver.onCompleted()
    }
}
