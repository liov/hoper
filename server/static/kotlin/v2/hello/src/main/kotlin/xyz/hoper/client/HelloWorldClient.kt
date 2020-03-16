package xyz.hoper.client

import io.grpc.ManagedChannel
import lombok.extern.log4j.Log4j2
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
import xyz.hoper.protobuf.GreeterGrpc
import xyz.hoper.protobuf.HelloReply
import xyz.hoper.protobuf.HelloRequest
import javax.annotation.PostConstruct

@Component
@Log4j2
class HelloWorldClient {
    @Value("\${grpc.client.host}")
    private val host: String = "localhost"

    @Value("\${grpc.client.port}")
    private val port: Int = 8080

    @Autowired
    private val grpcClientMananer: GrpcClientMananer? = null

    val log : Logger = LoggerFactory.getLogger(HelloWorldClient::class.java)

    @PostConstruct
    fun init() {
        call()
    }

    fun call() {
        val channel: ManagedChannel = grpcClientMananer!!.getChannel(host, port)
        val stub: GreeterGrpc.GreeterBlockingStub = GreeterGrpc.newBlockingStub(channel)
        val request: HelloRequest = HelloRequest.newBuilder().setName("world").build()
        val helloReply: HelloReply = stub.sayHello(request)
        log.info("time: " + helloReply.time)
    }
}
