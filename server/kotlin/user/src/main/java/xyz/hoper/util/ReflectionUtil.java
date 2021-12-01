package xyz.hoper.util;

import org.reflections.Reflections;
import org.reflections.util.ClasspathHelper;
import org.reflections.util.ConfigurationBuilder;
import org.reflections.util.FilterBuilder;

import java.util.Arrays;
import java.util.stream.Stream;

/**
 * reflections反射工
 */
public class ReflectionUtil {

    public static Reflections getReflections(String... packageAddress) {
        ConfigurationBuilder configurationBuilder = new ConfigurationBuilder();
        Stream.of(packageAddress).forEach(str -> configurationBuilder.addUrls(ClasspathHelper.forPackage(str.trim())));
        FilterBuilder filterBuilder = new FilterBuilder();
        Arrays.stream(packageAddress).forEach(filterBuilder::includePackage);
        configurationBuilder.filterInputsBy(filterBuilder);
        return new Reflections(configurationBuilder);
    }

    public static Reflections getReflections(String packageAddress) {
        ConfigurationBuilder configurationBuilder = new ConfigurationBuilder();
        Stream.of(packageAddress).forEach(str -> configurationBuilder.addUrls(ClasspathHelper.forPackage(str.trim())));
        FilterBuilder filterBuilder = new FilterBuilder();
        filterBuilder.includePackage(packageAddress);
        configurationBuilder.filterInputsBy(filterBuilder);
        return new Reflections(configurationBuilder);
    }

}
