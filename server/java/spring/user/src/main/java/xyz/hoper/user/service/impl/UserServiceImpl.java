package xyz.hoper.user.service.impl;


import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import xyz.hoper.user.dao.UserRepository;
import xyz.hoper.user.service.UserService;
import xyz.hoper.user.entity.User;


@Component
class UserServiceImpl implements UserService {

    @Autowired
    private UserRepository userRepository ;


    public User info(Long id ) {
        try {
            return userRepository.findById(id).get();
        } catch (Exception e) {
            e.printStackTrace();
        }
        return null;
    }
}
