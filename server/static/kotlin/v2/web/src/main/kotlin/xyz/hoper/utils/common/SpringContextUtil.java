package xyz.hoper.utils.common;

import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;
import org.springframework.stereotype.Component;

/**
 * Spring上下文工具
 */
@Component
public class SpringContextUtil implements ApplicationContextAware {

    private static ApplicationContext applicationContext;

    /**
     * 实现ApplicationContextAware接口的回调方法
     */
    @Override
    public void setApplicationContext(ApplicationContext applicationContext) throws BeansException {
        SpringContextUtil.applicationContext = applicationContext;
    }


    public static ApplicationContext getApplicationContext() {
        return applicationContext;
    }

    /**
     * 获取对象
     *
     * @return Object 一个以所给名字注册的bean的实例(必须遵循Spring的生成规则)
     */
    public static Object getBean(String name) throws BeansException {
        return applicationContext.getBean(name);
    }

    /**
     * 获取对象
     *
     * @return Object 一个以所给名字注册的bean的实例(必须遵循Spring的生成规则)
     */
    public static Object getBean(Class classObj) throws BeansException {
        // 获取bean的id
        // System.err.println("---->: "+ Arrays.toString(applicationContext.getBeanNamesForType(classObj)));
        return applicationContext.getBean(classObj);
    }
    
}
