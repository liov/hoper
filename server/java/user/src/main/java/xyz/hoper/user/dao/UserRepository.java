package xyz.hoper.user.dao;

import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;
import xyz.hoper.user.entity.User;


/**
 * @Description TODO
 * @Date 2022/11/21 10:43
 * @Created by lbyi
 */
@Repository
public interface UserRepository extends CrudRepository<User, Long> {}