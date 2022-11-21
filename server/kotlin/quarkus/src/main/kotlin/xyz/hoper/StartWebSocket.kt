package xyz.hoper


import jakarta.websocket.server.PathParam
import jakarta.websocket.server.ServerEndpoint

import javax.enterprise.context.ApplicationScoped
import javax.websocket.*

@ServerEndpoint("/start-websocket/{name}")
@ApplicationScoped
class StartWebSocket {

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
}
