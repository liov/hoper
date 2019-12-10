package xyz.hoper.web;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication(scanBasePackages={"xyz.hoper.web"})
@MapperScan("xyz.hoper.web.dao")
public class Application {
    /**
     *
     * @param args
     */
    public static void main(String[] args) {
        SpringApplication.run(Application.class, args);
    }
}
