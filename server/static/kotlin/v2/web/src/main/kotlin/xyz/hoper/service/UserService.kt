package xyz.hoper.service

import io.vertx.codegen.annotations.ProxyGen
import io.vertx.core.AsyncResult
import io.vertx.core.Handler
import xyz.hoper.entity.User


@ProxyGen
interface UserService {
    fun info(id: Long, resultHandler: Handler<AsyncResult<User>>)
}
