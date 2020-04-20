package xyz.hoper.api

import io.vertx.core.AsyncResult
import io.vertx.core.Handler
import io.vertx.ext.web.RoutingContext
import xyz.hoper.entity.User
import xyz.hoper.service.UserService
import xyz.hoper.utils.annotation.RouteHandler
import xyz.hoper.utils.annotation.RouteMapping
import xyz.hoper.utils.annotation.RouteMethod
import xyz.hoper.utils.resultvo.ResultBean
import xyz.hoper.utils.resultvo.ResultConstant
import xyz.hoper.utils.vertx.AsyncServiceUtil
import xyz.hoper.utils.vertx.HttpUtil
import java.net.HttpURLConnection.HTTP_INTERNAL_ERROR
import java.net.HttpURLConnection.HTTP_OK


@RouteHandler("user")
class UserApi {
    private val userAsyncService: UserService = AsyncServiceUtil.getAsyncServiceInstance(UserService::class.java)

    /**
     * 演示路径参数
     *
     * @return
     */
    @RouteMapping(value = "/:id", method = RouteMethod.GET)
    fun info(): Handler<RoutingContext> {
        return Handler<RoutingContext> { ctx ->
            userAsyncService.info(
                    ctx.get("id"), Handler { ar ->
                if (ar.succeeded()) {
                    val user = ar.result()
                    HttpUtil.jsonResponse(ctx.response(), HTTP_OK, ResultBean.build().setResultConstant(ResultConstant._000, user))
                } else {
                    HttpUtil.jsonResponse(ctx.response(), HTTP_INTERNAL_ERROR,
                            ResultBean.build().setData(ar.cause().message).setCode(HTTP_INTERNAL_ERROR.toString()))
                    ar.cause().printStackTrace()
                }
            })
        }
    }
}
