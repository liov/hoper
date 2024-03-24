package xyz.hoper.user.controller;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.*;
import reactor.core.publisher.Mono;
import xyz.hoper.user.entity.User;
import xyz.hoper.user.service.UserService;

/**
 * @Description TODO
 * @Date 2022/11/18 17:09
 * @Created by lbyi
 */
@RestController
@RequestMapping("/api/user")
class UserController {
    Logger log = LoggerFactory.getLogger(UserController.class);

    @Autowired
    private UserService userService;

    @GetMapping("{id}")
    User info(@PathVariable Long id) {
        return userService.info(id);
    }

}