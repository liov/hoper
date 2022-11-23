package xyz.hoper


import java.util.concurrent.ConcurrentHashMap
import javax.enterprise.context.ApplicationScoped
import javax.websocket.*
import javax.websocket.server.PathParam
import javax.websocket.server.ServerEndpoint

@ApplicationScoped
@ServerEndpoint("/start-websocket/{name}")
class StartWebSocket {

    var sessions: Map<String, Session> = ConcurrentHashMap()

    @OnOpen
    fun onOpen(session: Session?, @PathParam("name") name: String) {
        println("onOpen> $name")
    }

    @OnClose
    fun onClose(session: Session?, @PathParam("name") name: String) {
        println("onClose> $name")
    }

    @OnError
    fun onError(session: Session?, @PathParam("name") name: String, throwable: Throwable) {
        println("onError> $name: $throwable")
    }

    @OnMessage
    fun onMessage(message: String, @PathParam("name") name: String) {
        println("onMessage> $name: $message")
    }

    private fun broadcast(message: String) {
        sessions.values.forEach { s ->
            s.asyncRemote.sendObject(message) { result ->
                if (result.exception != null) {
                    println("Unable to send message: " + result.exception)
                }
            }
        }
    }
}
