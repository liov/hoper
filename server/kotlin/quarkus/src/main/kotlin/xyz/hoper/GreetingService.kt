package xyz.hoper

import io.quarkus.vertx.ConsumeEvent
import javax.enterprise.context.ApplicationScoped


/**
 * @Description TODO
 * @Date 2022/11/18 17:20
 * @Created by lbyi
 */
@ApplicationScoped
class GreetingService {
    @ConsumeEvent("greetings")
    fun hello(name: String): String {
        return "Hello $name"
    }
}