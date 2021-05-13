package xyz.hoper.user.service.impl

import io.vertx.core.AsyncResult
import io.vertx.core.Future
import io.vertx.core.Handler
import org.springframework.stereotype.Component
import xyz.hoper.web.entity.User
import xyz.hoper.web.service.UserService
import xyz.hoper.vertx.annotation.AsyncServiceHandler
import xyz.hoper.vertx.util.BaseAsyncService


@Component
@AsyncServiceHandler
class UserServiceImpl : UserService, BaseAsyncService {
    override fun info(id: Long, resultHandler: Handler<AsyncResult<User>>) {
        try {
            val user = User()
            user.id = id
            user.name = "测试"
            Future.succeededFuture(user).onComplete(resultHandler)
        } catch (e: Exception) {
            e.printStackTrace()
            resultHandler.handle(Future.failedFuture(e))
        }
    }
}
