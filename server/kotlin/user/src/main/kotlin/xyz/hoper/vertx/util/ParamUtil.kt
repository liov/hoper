package xyz.hoper.vertx.util

import io.vertx.core.MultiMap
import io.vertx.core.json.JsonArray
import io.vertx.core.json.JsonObject
import io.vertx.ext.web.RoutingContext
import org.slf4j.Logger
import org.slf4j.LoggerFactory

/**
 * 参数 工具类
 */
object ParamUtil {
    private val LOGGER: Logger = LoggerFactory.getLogger(ParamUtil::class.java)

    /**
     * 请求参数转换
     *
     * @param ctx 根据vertx-web的上下文来获取参数
     */
    fun getRequest(ctx: RoutingContext): JsonObject? {
        val paramMap: MultiMap = ctx.request().params()
        val param = JsonObject()
        //ip
        //params.put("serverIp", ctx.request().localAddress().host());
        //params.put("clientIp", ctx.request().remoteAddress().host());
        if (!paramMap.isEmpty) {
            val iter = paramMap.entries().iterator()
            while (iter.hasNext()) {
                val entry = iter.next() as Map.Entry<String, String>
                if (param.containsKey(entry.key)) { //多值
                    if (param.getValue(entry.key) !is JsonArray) {
                        param.put(entry.key, paramMap.getAll(entry.key))
                    }
                } else {
                    param.put(entry.key, entry.value)
                }
            }
        }
        if (param.isEmpty) LOGGER.debug("HttpServerRequest无请求参数! ")
        return getParamPage(param)
    }

    fun getRequestBody(ctx: RoutingContext): JsonObject {
        return ctx.body().asJsonObject()
    }

    /**
     * 默认处理分页
     *
     * @param params
     */
    private fun getParamPage(params: JsonObject?): JsonObject? {
        if (params != null) {
            if (!params.containsKey("limit")) {
                params.put("limit", 20)
            } else {
                var limit: Int = Integer.valueOf(params.getValue("limit").toString())
                limit = if (limit < 0) 20 else limit
                limit = if (limit > 50) 50 else limit
                params.put("limit", limit)
            }
            if (!params.containsKey("page")) {
                params.put("page", 1)
            } else {
                var page: Int = Integer.valueOf(params.getValue("page").toString())
                page = if (page < 0) 0 else page
                params.put("page", page)
            }
        }
        return params
    }
}