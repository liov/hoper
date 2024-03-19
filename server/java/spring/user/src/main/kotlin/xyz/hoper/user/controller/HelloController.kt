package xyz.hoper.user.controller

import org.slf4j.LoggerFactory
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RequestMapping
import org.springframework.web.bind.annotation.RestController
import reactor.core.publisher.Flux
import reactor.core.publisher.Mono
import reactor.core.publisher.MonoSink


@RestController
@RequestMapping("/api/hello")
class HelloController {
    val log = LoggerFactory.getLogger(this.javaClass)

    @GetMapping("mono")
    fun mono(): Mono<String> {
        return Mono.create { monoSink: MonoSink<String> ->
            log.info("创建 Mono")
            monoSink.success("hello webflux")
        }.doOnSubscribe { subscription ->
            log.info("{}", subscription)
        }.doOnNext { o -> log.info("{}", o) }
    }

    @GetMapping("flux")
    fun flux(): Flux<String> {
        return Flux.just("hello", "webflux", "spring", "boot")
    }
}