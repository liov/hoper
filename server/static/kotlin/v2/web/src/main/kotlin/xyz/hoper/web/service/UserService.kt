package xyz.hoper.web.service

import io.vertx.codegen.annotations.ProxyGen
import io.vertx.core.AsyncResult
import io.vertx.core.Handler
import xyz.hoper.web.entity.User


@ProxyGen
interface UserService {
    fun info(id: Long, resultHandler: Handler<AsyncResult<User>>)
}
