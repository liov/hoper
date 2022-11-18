package xyz.hoper.user.controller

import org.slf4j.Logger
import org.slf4j.LoggerFactory
import org.springframework.web.bind.annotation.*
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.core.publisher.MonoSink
import xyz.hoper.user.entity.User
import xyz.hoper.user.service.UserService

/**
 * @Description TODO
 * @Date 2022/11/18 17:09
 * @Created by lbyi
 */
@RestController
@RequestMapping("/api/user")
class UserApi {
    val log: Logger = LoggerFactory.getLogger(this.javaClass)
    private lateinit var webFluxService: UserService

    @GetMapping("{id}")
    fun mono(@PathVariable("id") id:Long): Mono<User> {
        return Mono.create { monoSink: MonoSink<User> ->
            log.info("创建 Mono")
            monoSink.success(webFluxService.info(id))
        }.doOnSubscribe { subscription ->
            log.info("{}", subscription)
        }.doOnNext { o -> log.info("{}", o) }
    }

}