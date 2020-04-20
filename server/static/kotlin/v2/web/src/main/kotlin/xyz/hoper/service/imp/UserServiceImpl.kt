package xyz.hoper.service.imp

import io.vertx.core.AsyncResult
import io.vertx.core.Future
import io.vertx.core.Handler
import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Component
import xyz.hoper.entity.User
import xyz.hoper.service.UserService
import xyz.hoper.utils.annotation.AsyncServiceHandler
import xyz.hoper.utils.vertx.BaseAsyncService


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
