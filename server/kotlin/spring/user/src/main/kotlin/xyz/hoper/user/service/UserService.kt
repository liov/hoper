package xyz.hoper.user.service


import reactor.core.publisher.Mono
import xyz.hoper.user.entity.User


interface UserService {
    fun info(id: Long): Mono<User>
}
