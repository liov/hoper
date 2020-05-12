package xyz.hoper.unsafe;

import sun.misc.Unsafe;

import java.lang.reflect.Field;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 9:27
 * @description：unsafe
 * @modified By：
 */
public class MyObj {
    int objField = 10;
    static int clsField = 10;
    int[] array = {10, 20, 30, 40, 50};
    static Unsafe U;

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

    static void getObjFieldVal() throws NoSuchFieldException {
        Field field = MyObj.class.getDeclaredField("objField");
        long offset = U.objectFieldOffset(field);
        MyObj obj = new MyObj();

        int val = U.getInt(obj, offset);
        System.out.println("1.\t" + val + "\t" + (val == 10));

        MyObjChild child = new MyObjChild();
        int corVal1 = U.getInt(child, offset);
        System.out.println("2.\t" + corVal1 + "\t" + (corVal1 == 10));

        Field fieldChild = MyObj.class.getDeclaredField("objField");
        long offsetChild = U.objectFieldOffset(fieldChild);
        System.out.println("3.\t" + offset + "\t" + offsetChild + "\t" + (offset == offsetChild));
        int corVal2 = U.getInt(obj, offsetChild);
        System.out.println("4.\t" + corVal2 + "\t" + (corVal2 == 10));

        short errVal1 = U.getShort(obj, offset);
        System.out.println("5.\t" + errVal1 + "\t" + (errVal1 == 10));

        int errVal2 = U.getInt("abcd", offset);
        System.out.println("6.\t" + errVal2);
    }

    static void getClsFieldVal() throws NoSuchFieldException {
        Field field = MyObj.class.getDeclaredField("clsField");
        long offset = U.staticFieldOffset(field);
        Object obj = U.staticFieldBase(field);
        int val1 = U.getInt(MyObj.class, offset);
        System.out.println("1.\t" + val1 + "\t" + (val1 == 10));
        int val2 = U.getInt(obj, offset);
        System.out.println("2.\t" + val2 + "\t" + (val2 == 10));
    }

    static void getArrayVal(int index, int expectedVal) throws NoSuchFieldException {
        int base = U.arrayBaseOffset(int[].class);
        int scale = U.arrayIndexScale(int[].class);
        MyObj obj = new MyObj();
        Field field = MyObj.class.getDeclaredField("array");
        long offset = U.objectFieldOffset(field);
        int[] array = (int[]) U.getObject(obj, offset);
        int val1 = U.getInt(array, (long) base + index * scale);
        System.out.println("1.\t" + val1 + "\t" + (val1 == expectedVal));
        int val2 = U.getInt(obj.array, (long) base + index * scale);
        System.out.println("2.\t" + val2 + "\t" + (val2 == expectedVal));
    }

    static void setObjFieldVal(int val) throws NoSuchFieldException {
        Field field = MyObj.class.getDeclaredField("objField");
        long offset = U.objectFieldOffset(field);
        MyObj obj = new MyObj();
        U.putInt(obj, offset, val);
        int getVal = U.getInt(obj, offset);
        System.out.println(val + "\t" + (val == getVal));
        //U.putLong(obj, offset, val);
        /**
         * 不注释后面代码会报错，因为修改到了array的内存值。
         * objField 是 int 类型，通过 putInt 修改 objField 的值可以正常修改，结果打印 true。
         * 然后使用 putLong 修改 objField 的值，这个修改操作本身也不会报错，
         * 但是 objField 并不是 long 类型，这样修改会导致其它程序错误，
         * 它不仅修改了 objField 的内存值，还修改了 objField 之后四个字节的内存值。
         * 在这个例子中，objField 后面八个字节存储的是 array 字段所表示的数组对象的偏移位置，但是被修改了，
         * 如果后面的代码尝试访问 array 字段就会出错。修改其它字段（array、clsField）也是一样的，
         * 只要按之前的方法获取字段的偏移位置，使用与字段类型一致的 put 方法就可以。
         */
        Field fieldArray = MyObj.class.getDeclaredField("array");
        long offsetArray = U.objectFieldOffset(fieldArray);
        int[] array = (int[]) U.getObject(obj, offsetArray);
        for (int i = 0; i < array.length; i++) {
            System.out.println(array[i]);
        }
    }

    static void memory() {
        long address = U.allocateMemory(10);
        U.setMemory(address, 10, (byte) 1);
        /**
         * 1的二进制码为00000001，int为四个字节，U.getInt将读取四个字节，
         * 读取的字节为00000001 00000001 00000001 00000001
         */
        int i = 0b00000001000000010000000100000001;
        System.out.println(i + "\t" + U.getInt(address));
        U.freeMemory(address);
    }

    static void clsInitialized() throws NoSuchFieldException {
        System.out.println(U.shouldBeInitialized(MyObj.class));
        System.out.println(U.shouldBeInitialized(MyObjChild.class));
        U.ensureClassInitialized(MyObjChild.class);
        System.out.println(U.shouldBeInitialized(MyObjChild.class));
    }

    static void clsInitialized2() throws NoSuchFieldException {
        Field f1=MyObjChild.class.getDeclaredField("f1");
        Field f2=MyObjChild.class.getDeclaredField("f2");
        long f1Offset= U.staticFieldOffset(f1);
        long f2Offset= U.staticFieldOffset(f2);
        int f1Val=U.getInt(MyObjChild.class,f1Offset);
        int f2Val=U.getInt(MyObjChild.class,f2Offset);
        System.out.println("1.\t"+(f1Val==0));
        System.out.println("2.\t"+(f2Val==1));
        U.ensureClassInitialized(MyObjChild.class);
        f1Val=U.getInt(MyObjChild.class,f1Offset);
        System.out.println("3.\t"+(f1Val==2));
    }

    static void cas() throws NoSuchFieldException {
        Field field=MyObj.class.getDeclaredField("objField");
        long offset= U.objectFieldOffset(field);
        MyObj myObj=new MyObj();
        myObj.objField=1;
        U.compareAndSwapInt(myObj,offset,0,2);
        System.out.println("1.\t"+(myObj.objField==2));
        U.compareAndSwapInt(myObj,offset,1,2);
        System.out.println("2.\t"+(myObj.objField==2));
    }

    static void clsInitialized3() throws InstantiationException {
        MyObjChild myObj= (MyObjChild) U.allocateInstance(MyObjChild.class);
        System.out.println("1.\t"+(MyObjChild.f1==2));
        System.out.println("1.\t"+(myObj.f2==0));
    }
}
