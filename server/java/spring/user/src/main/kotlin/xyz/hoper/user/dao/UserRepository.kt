package xyz.hoper.user.dao

import org.springframework.data.repository.reactive.ReactiveCrudRepository
import org.springframework.stereotype.Repository
import xyz.hoper.user.entity.User


/**
 * @Description TODO
 * @Date 2022/11/21 10:43
 * @Created by lbyi
 */
@Repository
interface UserRepository : ReactiveCrudRepository<User, Long>