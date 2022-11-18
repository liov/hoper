package xyz.hoper.user


import org.springframework.boot.runApplication
import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.context.annotation.ComponentScan

@ComponentScan("xyz.hoper.util","xyz.hoper.user")
@SpringBootApplication
class UserApplication {
}


fun main(args: Array<String>) {
    runApplication<UserApplication>(*args)
}
