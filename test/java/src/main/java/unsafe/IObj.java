package unsafe;

import jdk.internal.misc.Unsafe;

import java.lang.reflect.Field;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 14:43
 * @description：
 * @modified By：
 */
public class IObj {
    private int objField = 10;
    private static Unsafe U = Unsafe.getUnsafe();

    //java.exe --add-exports=java.base/jdk.internal.misc=ALL-UNNAMED IObj.java
    //IDEA中VM options：--add-exports=java.base/jdk.internal.misc=test
    public static void main(String[] args) throws NoSuchFieldException {
        Field field = IObj.class.getDeclaredField("objField");
        long offset = U.objectFieldOffset(field);
        IObj obj = new IObj();
        int val = U.getInt(obj, offset);
        System.out.println("1.\t" + val + "\t" + (val == 10));
        U.putInt(obj, offset,11);
        System.out.println(U.getAddress(obj,offset));
        System.out.println("2.\t" + U.getInt(obj, offset) + "\t" + (U.getInt(obj, offset) == 10));
    }
}
