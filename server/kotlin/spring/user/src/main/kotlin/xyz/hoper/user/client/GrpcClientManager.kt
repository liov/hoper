package xyz.hoper.user.client


import java.util.concurrent.TimeUnit
import io.grpc.ManagedChannel
import io.grpc.ManagedChannelBuilder
import org.springframework.stereotype.Component

@Component
class GrpcClientManager {
    fun getChannel(host: String?, port: Int): ManagedChannel? {
        return ManagedChannelBuilder.forAddress(host, port)
                .usePlaintext()
                .disableRetry()
                .idleTimeout(5, TimeUnit.SECONDS)
                .build()
    }
}