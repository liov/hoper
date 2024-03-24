package xyz.hoper.user.config;

import org.springframework.boot.context.properties.ConfigurationProperties;
import org.springframework.stereotype.Component;

/**
 * @Description TODO
 * @Date 2022/11/21 15:29
 * @Created by lbyi
 */

@ConfigurationProperties(prefix = "config")
@Component
class UserProperties {
    String name;
    String password;
    String ip;
    Integer port;
}
