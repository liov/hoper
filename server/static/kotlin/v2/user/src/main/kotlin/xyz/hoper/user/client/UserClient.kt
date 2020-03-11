package xyz.hoper.user.client


import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.beans.factory.annotation.Value
import org.springframework.stereotype.Component
import javax.annotation.PostConstruct
import xyz.hoper.protobuf.user.UserServiceGrpc
import xyz.hoper.protobuf.user.UserServiceOuterClass


@Component
class UserClient {
     @Value("\${grpc.client.host}")
     private var host: String? = null
    @Value("\${grpc.client.port}")
    private var port: Int = 8080

    @Autowired
    private var grpcClientMananer: GrpcClientMananer? = null

    val log : Logger = LoggerFactory.getLogger(UserClient::class.java)

    @PostConstruct
    fun init() {
        call()
    }

    fun call(){
        val channel = grpcClientMananer?.getChannel(host,port);
        val stub = UserServiceGrpc.newBlockingStub(channel);
        val request = UserServiceOuterClass.GetReq.newBuilder().setId(1).build();
        val reply = stub.getUser(request);
        log.info("time: "+ reply);
    }

}