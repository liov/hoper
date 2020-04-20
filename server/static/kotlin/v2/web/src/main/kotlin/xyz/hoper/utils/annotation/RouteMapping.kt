package xyz.hoper.utils.annotation


/**
 * Router API Mehtod 标识注解
 */
@Target(AnnotationTarget.FUNCTION, AnnotationTarget.PROPERTY_GETTER, AnnotationTarget.PROPERTY_SETTER)
@Retention(AnnotationRetention.RUNTIME)
@MustBeDocumented
annotation class RouteMapping(
        val value: String = "",
        /**** 是否覆盖  */
        val isCover: Boolean = true,
        /**** 使用http method  */
        val method: RouteMethod = RouteMethod.GET,
        /**** 接口描述  */
        val descript: String = "",
        /**
         * 注册顺序，数字越大越先注册
         */
        val order: Int = 0
)