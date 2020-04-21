package xyz.hoper.web.api

import io.vertx.core.Handler
import io.vertx.core.json.JsonObject
import io.vertx.ext.web.RoutingContext
import xyz.hoper.vertx.annotation.RouteHandler
import xyz.hoper.vertx.annotation.RouteMapping
import xyz.hoper.vertx.annotation.RouteMethod
import xyz.hoper.vertx.resultvo.ResultBean
import xyz.hoper.vertx.resultvo.ResultConstant
import xyz.hoper.vertx.util.AsyncServiceUtil
import xyz.hoper.vertx.util.HttpUtil
import xyz.hoper.vertx.util.ParamUtil
import xyz.hoper.web.service.UserService
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
            val param: JsonObject = ParamUtil.getRequest(ctx)!!
            if (!param.containsKey("id")) {
                HttpUtil.jsonResponse(ctx.response(), HTTP_INTERNAL_ERROR,
                        ResultBean.build().setMsg("缺少id参数").setCode(HTTP_INTERNAL_ERROR.toString()));
                return@Handler
            }
            userAsyncService.info(
                   param.getString("id").toLong(), Handler { ar ->
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

    @RouteMapping(value = "/:id", method = RouteMethod.GET)
    fun test(): Handler<RoutingContext>? {
        return Handler { ctx: RoutingContext ->
            val param: JsonObject? = ParamUtil.getRequest(ctx)
            if (param != null) {
                HttpUtil.jsonResponse(ctx.response(), HTTP_OK, ResultBean.build().setMsg("Hello，欢迎使用测试地址.....").setData(param.encode()))
            }
        }
    }
}
