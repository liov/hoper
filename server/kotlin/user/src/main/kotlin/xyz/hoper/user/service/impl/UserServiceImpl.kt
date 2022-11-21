package xyz.hoper.user.service.impl


import org.springframework.beans.factory.annotation.Autowired
import org.springframework.stereotype.Component
import reactor.core.publisher.Mono
import xyz.hoper.user.dao.UserRepository
import xyz.hoper.user.service.UserService
import xyz.hoper.user.entity.User


@Component
class UserServiceImpl : UserService {

    @Autowired
    private lateinit var userRepository: UserRepository
    override fun info(id: Long): Mono<User> {
        try {
            return userRepository.findById(id)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return Mono.empty()
    }
}
