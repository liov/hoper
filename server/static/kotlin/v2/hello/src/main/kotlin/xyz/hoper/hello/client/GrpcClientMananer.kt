package xyz.hoper.hello.client

import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import org.springframework.stereotype.Component
import java.util.concurrent.TimeUnit

@Component
class GrpcClientMananer {
    fun getChannel(host: String?, port: Int): ManagedChannel {
        return ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext()
                .disableRetry()
                .idleTimeout(5, TimeUnit.SECONDS)
                .build()
    }
}