package xyz.hoper.user.service


import xyz.hoper.user.entity.User



interface UserService {
    fun info(id: Long):User?
}
