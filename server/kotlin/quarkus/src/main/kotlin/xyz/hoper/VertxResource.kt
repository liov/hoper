package xyz.hoper

import io.smallrye.mutiny.Multi
import io.smallrye.mutiny.Uni
import io.vertx.core.file.OpenOptions
import io.vertx.core.json.JsonArray
import io.vertx.core.json.JsonObject

import io.vertx.mutiny.core.Vertx
import io.vertx.mutiny.core.buffer.Buffer
import io.vertx.mutiny.core.eventbus.EventBus
import io.vertx.mutiny.core.file.AsyncFile
import io.vertx.mutiny.ext.web.client.HttpResponse
import io.vertx.mutiny.ext.web.client.WebClient
import java.nio.charset.StandardCharsets
import javax.inject.Inject
import javax.ws.rs.GET
import javax.ws.rs.Path
import javax.ws.rs.QueryParam


/**
 * @Description TODO
 * @Date 2022/11/18 16:53
 * @Created by lbyi
 */

@Path("/vertx")
class VertxResource @Inject constructor(private val vertx: Vertx) {

    private val client = WebClient.create(vertx)

    @GET
    @Path("/lorem")
    fun readShortFile(): Uni<String> {
        return vertx.fileSystem().readFile("lorem.txt")
            .onItem().transform { content: Buffer ->
                content.toString(
                    StandardCharsets.UTF_8
                )
            }
    }

    @GET
    @Path("/book")
    fun readLargeFile(): Multi<String> {
        return vertx.fileSystem().open(
            "book.txt",
            OpenOptions().setRead(true)
        )
            .onItem().transformToMulti { file: AsyncFile -> file.toMulti() }
            .onItem().transform { content: Buffer ->
                """
                ${content.toString(StandardCharsets.UTF_8)}
                ------------
                
                """.trimIndent()
            }
    }

    @Inject
    lateinit var bus: EventBus

    @GET
    @Path("/hello")
    fun hello(@QueryParam("name") name: String?): Uni<String> {
        return bus.request<String>("greetings", name).onItem().transform { response -> response.body() }
    }

    private val URL = "https://en.wikipedia.org/w/api.php?action=parse&page=Quarkus&format=json&prop=langlinks"

    @GET
    @Path("/web")
    fun retrieveDataFromWikipedia(): Uni<JsonArray> {
        return client.getAbs(URL).send()
            .onItem().transform(HttpResponse<Buffer>::bodyAsJsonObject)
            .onItem().transform { json: JsonObject ->
                json.getJsonObject("parse")
                    .getJsonArray("langlinks")
            }
    }

}