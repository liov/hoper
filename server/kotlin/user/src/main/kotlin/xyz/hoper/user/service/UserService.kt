package xyz.hoper.user.service


import org.springframework.stereotype.Service
import reactor.core.publisher.Mono
import xyz.hoper.user.entity.User


interface UserService {
    fun info(id: Long): Mono<User>
}
