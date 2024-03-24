package xyz.hoper.user.service;


import reactor.core.publisher.Mono;
import xyz.hoper.user.entity.User;


public interface UserService {
     User info(Long id);
}
