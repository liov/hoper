package xyz.hoper.user.api

import io.vertx.core.Handler
import io.vertx.core.file.FileSystem
import io.vertx.ext.web.FileUpload
import io.vertx.ext.web.RoutingContext
import xyz.hoper.vertx.annotation.RouteHandler
import xyz.hoper.vertx.annotation.RouteMapping
import xyz.hoper.vertx.annotation.RouteMethod
import xyz.hoper.vertx.resultvo.ResultBean
import xyz.hoper.vertx.resultvo.ResultConstant
import xyz.hoper.vertx.util.HttpUtil
import xyz.hoper.vertx.util.VertxUtil.vertxInstance
import java.net.HttpURLConnection.HTTP_OK
import java.util.function.Consumer


@RouteHandler
class UploadAndDownApi {
    /**
     * 文件上传
     *
     * @return
     */
    @RouteMapping(value = "/upload", method = RouteMethod.POST)
    fun upload(): Handler<RoutingContext> {
        return Handler { ctx ->
            val uploads: List<FileUpload> = ctx.fileUploads()
            val fs: FileSystem = vertxInstance!!.fileSystem()
            uploads.forEach(Consumer<FileUpload> { fileUpload: FileUpload ->
                val path = "/upload/" + fileUpload.fileName()
                fs.copy(fileUpload.uploadedFileName(), path) { ar ->
                    if (ar.succeeded()) {
                        fs.deleteBlocking(fileUpload.uploadedFileName())
                    }
                }
            })
            HttpUtil.jsonResponse(ctx.response(), HTTP_OK, ResultBean.build().setResultConstant(ResultConstant.SUCCESS))
        }
    }

    /**
     * 文件下载
     * @return
     */
    @RouteMapping(value = "/down", method = RouteMethod.GET)
    fun download(): Handler<RoutingContext> {
        return Handler { ctx ->
            ctx.response().putHeader("content-Type", "application/x-png")
            ctx.response().putHeader("Content-Disposition", "attachment;filename=" + "hahaha.png")
            ctx.response().sendFile(ctx.queryParams().get("filepath"))
        }
    }
}
