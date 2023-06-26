package xyz.hoper.user.client


import io.grpc.StatusRuntimeException
import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
import xyz.hoper.pandora.protobuf.empty.EmptyOuterClass
import javax.annotation.PostConstruct
import xyz.hoper.protobuf.user.UserServiceGrpc


@Component
class UserClient {
     @Value("\${grpc.client.host}")
     private var host: String? = null
    @Value("\${grpc.client.port}")
    private var port: Int = 8080

    @Autowired
    private var grpcClientManager: GrpcClientManager? = null

    val log : Logger = LoggerFactory.getLogger(UserClient::class.java)

    @PostConstruct
    fun init() {
        call()
    }

    fun call(){
        val channel = grpcClientManager?.getChannel(host,port)
        val stub = UserServiceGrpc.newBlockingStub(channel)
        val request = EmptyOuterClass.Empty.newBuilder().build()
        try{
            val reply = stub.verifyCode(request)
            log.info("time: $reply")
        }catch (e: StatusRuntimeException) {
            log.error(e.message)
        }
    }

}