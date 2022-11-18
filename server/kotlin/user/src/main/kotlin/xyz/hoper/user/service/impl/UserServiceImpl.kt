package xyz.hoper.user.service.impl


import org.springframework.stereotype.Component
import xyz.hoper.user.service.UserService
import xyz.hoper.user.entity.User


@Component
class UserServiceImpl : UserService {
    override fun info(id: Long): User? {
        try {
            val user = User()
            user.id = id
            user.name = "测试"
            return user
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return null
    }
}
