package xyz.hoper.vertx.util

import io.vertx.core.http.HttpServerResponse
import io.vertx.ext.web.RoutingContext
import xyz.hoper.vertx.resultvo.ResultBean

object HttpUtil {
    /**
     * json格式
     */
    fun jsonResponse(response: HttpServerResponse, resultBean: ResultBean) {
        response.putHeader("content-type", "application/json; charset=utf-8").end(resultBean.toString())
    }

    fun jsonResponse(response: HttpServerResponse, statusCode: Int, resultBean: ResultBean) {
        response.putHeader("content-type", "application/json; charset=utf-8").setStatusCode(statusCode).end(resultBean.toString())
    }

    /**
     * 文本数据String
     */
    fun textResponse(routingContext: RoutingContext, text: String?) {
        routingContext.response().putHeader("content-type", "text/html; charset=utf-8").end(text)
    }
}