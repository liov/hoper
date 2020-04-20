package xyz.hoper.utils.common

import org.reflections.Reflections
import org.reflections.util.ClasspathHelper
import org.reflections.util.ConfigurationBuilder
import org.reflections.util.FilterBuilder
import java.util.stream.Stream

/**
 * reflections反射工具
 */
object ReflectionUtil {
    fun getReflections(vararg packageAddress: String?): Reflections {
        val configurationBuilder = ConfigurationBuilder()
        val filterBuilder = FilterBuilder()
        packageAddress.forEach { str -> configurationBuilder.addUrls(ClasspathHelper.forPackage(str?.trim { it <= ' ' })) }
        filterBuilder.includePackage(*packageAddress)
        configurationBuilder.filterInputsBy(filterBuilder)
        return Reflections(configurationBuilder)
    }
}