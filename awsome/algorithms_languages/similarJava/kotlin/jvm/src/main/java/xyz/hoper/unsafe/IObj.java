package xyz.hoper.unsafe;

import sun.misc.Unsafe;

import java.lang.reflect.Field;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 14:43
 * @description ：sun.misc.Unsafe用的就是jdk.internal.misc.Unsage,
 * 模块化jdk.internal.misc.Unsage无法导出，只能用sun.misc.Unsafe了
 * @modified By：
 */
public class IObj {
    private int objField = 10;
    private static Unsafe U;

    //源码getUnsafe方法有这个!VM.isSystemDomainLoader(caller.getClassLoader())
    //判断，必须要由启动类加载器加载,所以反射获取

    /**
     * 启动（Bootstrap）类加载器：启动类加载器是用本地代码实现的类加载器，它负责将JAVA_HOME/lib下面的核心类库或-Xbootclasspath选项指定的jar包等虚拟机识别的类库加载到内存中。由于启动类加载器涉及到虚拟机本地实现细节，开发者无法直接获取到启动类加载器的引用。具体可由启动类加载器加载到的路径可通过System.getProperty(“sun.boot.class.path”)查看。
     * 扩展（Extension）类加载器：扩展类加载器是由Sun的ExtClassLoader（sun.misc.Launcher$ExtClassLoader）实现的，它负责将JAVA_HOME /lib/ext或者由系统变量-Djava.ext.dir指定位置中的类库加载到内存中。开发者可以直接使用标准扩展类加载器，具体可由扩展类加载器加载到的路径可通过System.getProperty("java.ext.dirs")查看。
     * 系统（System）类加载器：系统类加载器是由 Sun 的 AppClassLoader（sun.misc.Launcher$AppClassLoader）实现的，它负责将用户类路径(java -classpath或-Djava.class.path变量所指的目录，即当前类所在路径及其引用的第三方类库的路径，如第四节中的问题6所述)下的类库加载到内存中。开发者可以直接使用系统类加载器，具体可由系统类加载器加载到的路径可通过System.getProperty("java.class.path")查看。
     * 线程上下文类加载器（context class loader）是从 JDK 1.2 开始引入的。Java.lang.Thread中的方法 getContextClassLoader()和 setContextClassLoader(ClassLoader cl)用来获取和设置线程的上下文类加载器。如果没有通过 setContextClassLoader(ClassLoader cl)方法进行设置的话，线程将继承其父线程的上下文类加载器。Java 应用运行的初始线程的上下文类加载器是系统类加载器，在线程中运行的代码可以通过此类加载器来加载类和资源。
     *
     */
    static {
        try {
            init();
        } catch (NoSuchFieldException | IllegalAccessException e) {
            e.printStackTrace();
        }
    }

    public static void init() throws NoSuchFieldException, IllegalAccessException {
        Field f = Unsafe.class.getDeclaredField("theUnsafe");
        f.setAccessible(true);
        U = (Unsafe) f.get(null);
    }

    //java.exe --add-exports=java.base/jdk.internal.misc=ALL-UNNAMED IObj.java
    //IDEA中VM options：--add-exports=java.base/jdk.internal.misc=test
    public static void main(String[] args) throws NoSuchFieldException {
        Field field = IObj.class.getDeclaredField("objField");
        long offset = U.objectFieldOffset(field);
        IObj obj = new IObj();
        int val = U.getInt(obj, offset);
        System.out.println("1.\t" + val + "\t" + (val == 10));
        U.putInt(obj, offset,11);
        //System.out.println(U.getAddress(obj, offset));//不开放，没懂理由
        System.out.println("2.\t" + U.getInt(obj, offset) + "\t" + (U.getInt(obj, offset) == 10));
    }
}
