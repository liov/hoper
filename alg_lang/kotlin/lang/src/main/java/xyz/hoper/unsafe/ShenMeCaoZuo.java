package xyz.hoper.unsafe;

import org.objenesis.Objenesis;
import org.objenesis.ObjenesisStd;
import sun.misc.Unsafe;

import java.lang.reflect.InvocationTargetException;
import java.lang.reflect.Method;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 11:29
 * @description：
 * @modified By：
 */
public class ShenMeCaoZuo {

    boolean override;

    private static final Objenesis OBJENESIS = new ObjenesisStd();

    public static void main(String[] args) throws ClassNotFoundException, NoSuchMethodException, NoSuchFieldException, InvocationTargetException, IllegalAccessException {
        Class<?> unsafeClass = Class.forName("jdk.internal.misc.Unsafe");
        Object unsafe = OBJENESIS.newInstance(unsafeClass);
        Unsafe sunMiscUnsafe = OBJENESIS.newInstance(sun.misc.Unsafe.class);
        Method unalignedAccess = unsafeClass.getDeclaredMethod("unalignedAccess");
        /**
         * jdk11正常,12报错Exception in thread "main" java.lang.NoSuchFieldException: override
         * long overrideOffset = sunMiscUnsafe.objectFieldOffset(AccessibleObject.class.getDeclaredField("override"));
         * NOTE: for security purposes, this field must not be visible outside this package.
         */
        long overrideOffset = sunMiscUnsafe.objectFieldOffset(ShenMeCaoZuo.class.getDeclaredField("override"));
        sunMiscUnsafe.putBoolean(
                unsafeClass,
                overrideOffset,
                true);
        sunMiscUnsafe.putBoolean(
                unalignedAccess,
                overrideOffset,
                true);
        System.out.println(unalignedAccess.invoke(unsafe));
    }
}
