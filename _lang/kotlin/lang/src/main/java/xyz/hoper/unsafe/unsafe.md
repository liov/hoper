
JDK Unsafe 源码完全注释 原
 编辑部的故事   编辑部的故事 发布于 03/08 15:19 字数 9065 阅读 3422 收藏 183 点赞 18  评论 4
JavaJDKCAS


并发作为 Java 中非常重要的一部分，其内部大量使用了 Unsafe 类，它为 java.util.concurrent 包中的类提供了底层支持。然而 Unsafe 并不是 JDK 的标准，它是 Sun 的内部实现，存在于 sun.misc 包中，在 Oracle 发行的 JDK 中并不包含其源代码。

Unsafe 提供两个功能：

绕过 JVM 直接修改内存（对象）
使用硬件 CPU 指令实现 CAS 原子操作


虽然我们在一般的并发编程中不会直接用到 Unsafe，但是很多 Java 基础类库与诸如 Netty、Cassandra 和 Kafka 等高性能库都采用它，它在提升 Java 运行效率、增强 Java 语言底层操作能力方面起了很大作用。笔者觉得了解一个使用如此广泛的库还是很有必要的。本文将深入到 Unsafe 的源码，分析一下它的逻辑。

本文使用 OpenJDK（jdk8-b120）中 Unsafe 的源码，Unsafe 的实现是和虚拟机实现相关的，不同的虚拟机实现，它们的对象结构可能不一样，这个 Unsafe 只能用于 Hotspot 虚拟机。

源码查看：http://hg.openjdk.java.net/jdk/jdk/file/a1ee9743f4ee/jdk/src/share/classes/sun/misc/Unsafe.java

上源码
Unsafe 为调用者提供执行非安全操作的能力，由于返回的 Unsafe 对象可以读写任意的内存地址数据，调用者应该小心谨慎的使用改对象，一定不用把它传递到非信任代码。该类的大部分方法都是非常底层的操作，并牵涉到一小部分典型的机器都包含的硬件指令，编译器可以对这些进行优化。

    public final class Unsafe {
    
        private static native void registerNatives();
        static {
            registerNatives();
            sun.reflect.Reflection.registerMethodsToFilter(Unsafe.class, "getUnsafe");
        }
    
        private Unsafe() {}
    
        private static final Unsafe theUnsafe = new Unsafe();
    
    
        @CallerSensitive
        public static Unsafe getUnsafe() {
            Class<?> caller = Reflection.getCallerClass();
            if (!VM.isSystemDomainLoader(caller.getClassLoader()))
                throw new SecurityException("Unsafe");
            return theUnsafe;
        }
        ......
    }
上面的代码包含如下功能：

本地静态方法：registerNatives()，该方法会在静态块中执行
私有构造函数：该类实例是单例的，不能实例化，可以通过 getUnsafe() 方法获取实例
静态单例方法：getUnsafe()，获取实例
静态块：包含初始化的注册功能
要使用此类必须获得其实例，获得实例的方法是 getUnsafe()，那就先看看这个方法。

    getUnsafe() 方法包含一个注释 @CallerSensitive，说明该方法不是谁都可以调用的。如果调用者不是由系统类加载器（bootstrap classloader）加载，则将抛出 SecurityException，所以默认情况下，应用代码调用此方法将抛出异常。我们的代码要想通过 getUnsafe() 获取实例是不可能的了，不过可通过反射获取 Unsafe 实例：
    
    Field f= Unsafe.class.getDeclaredField("theUnsafe");
    f.setAccessible(true);
    U= (Unsafe) f.get(null);
    
此处通过反射获取类的静态字段，这样就绕过了 getUnsafe() 的安全限制。

也可以通过反射获取构造方法再实例化，但这样违法了该类单例的原则，并且在使用上可能会有其它问题，所以不建议这样做。

 

再来看看如何获取指定变量的值：

public native int getInt(Object o, long offset);
获取指定对象中指定偏移量的字段或数组元素的值，参数 o 是变量关联的 Java 堆对象，如果指定的对象为 null，则获取内存中该地址的值（即 offset 为内存绝对地址）。

如果不符合以下任意条件，则该方法返回的值是不确定的：

offset 通过 objectFieldOffset 方法获取的类的某一字段的偏移量，并且关联的对象 o 是该类的兼容对象（对象 o 所在的类必须是该类或该类的子类）
offset 通过 staticFieldOffset 方法获取的类的某一字段的偏移量，o 是通过 staticFieldBase 方法获取的对象
如果 o 引用的是数组，则 offset 的值为 test.init.B+N*S，其中 N 是数组的合法下标，test.init.B 是通过 arrayBaseOffset 方法从该数组类获取的，S 是通过 arrayIndexScale 方法从该数组类获取的
如果以上任意条件符合，则调用者能获取 Java 字段的引用，但是如果该字段的类型和该方法返回的类型不一致，则结果是不一定的，比如该字段是 short，但调用了 getInt 方法。

该方法通过两个参数引用一个变量，它为 Java 变量提供 double-register 地址模型。如果引用的对象为 null，则该方法将 offset 当作内存绝对地址，就像 getInt(long)一样，它为非 Java 变量提供 single-register 地址模型，然而 Java 变量的内存布局可能和非 Java 变量的内存布局不同，不应该假设这两种地址模型是相等的。同时，应该记住 double-register 地址模型的偏移量不应该和 single-register 地址模型中的地址（long 参数）混淆。

 

再看条件中提到的几个相关方法：

    public native long objectFieldOffset(Field f);
    
    public native Object staticFieldBase(Field f);
    
    public native long staticFieldOffset(Field f);
    
    public native int arrayBaseOffset(Class arrayClass);
    
    public native int arrayIndexScale(Class arrayClass);
这几个方法分别是获取静态字段、非静态字段与数组字段的一些信息。

objectFieldOffset

很难想象 JVM 需要使用这么多比特位来编码非数组对象的偏移量，为了和该类的其它方法保持一致，所以该方法也返回 long 类型。

staticFieldBase

获取指定静态字段的位置，和 staticFieldOffset 一起使用。获取该静态字段所在的“对象”可通过类似 getInt(Object,long)的方法访问。

staticFieldOffset

返回给定的字段在该类的偏移地址。对于任何给定的字段，该方法总是返回相同的值，同一个类的不同字段总是返回不同的值。从 1.4.1 开始，字段的偏移以 long 表示，虽然 Sun 的 JVM 只使用了 32 位，但是那些将静态字段存储到绝对地址的 JVM 实现需要使用 long 类型的偏移量，通过 getXX(null,long) 获取字段值，为了保持代码迁移到 64 位平台上 JVM 的优良性，必须保持静态字段偏移量的所有比。

arrayBaseOffset

返回给定数组类第一个元素在内存中的偏移量。如果 arrayIndexScale 方法返回非 0 值，要获得访问数组元素的新的偏移量，则需要使用 s。

arrayIndexScale

返回给定数组类的每个元素在内存中的 scale(所占用的字节)。然而对于“narrow”类型的数组，类似 getByte(Object, int)的访问方法一般不会获得正确的结果，所以这些类返回的 scale 会是 0。

下边用代码解释：

    public class MyObj {
        int objField=10;
        static int clsField=10;
        int[] array={10,20,30,40,50};
        static Unsafe U;
        static {
            try {
                init();
            } catch (NoSuchFieldException e) {
                e.printStackTrace();
            } catch (IllegalAccessException e) {
                e.printStackTrace();
            }
        }
        public static void init() throws NoSuchFieldException, IllegalAccessException {
            Field f= Unsafe.class.getDeclaredField("theUnsafe");
            f.setAccessible(true);
            U= (Unsafe) f.get(null);
        }
        ......
    }
定义一个类包含成员变量 objField、类变量 clsField、成员数组 array 用于实验。要获取正确的结果，必须满足注释里的三个条件之一：

1、offset 通过 objectFieldOffset 方法获取的类的某一字段的偏移量，并且关联的对象 o 是该类的兼容对象(对象 o 所在的类必须是该类或该类的子类)

public class MyObjChild extends MyObj {
    int anthor;
}
    static void getObjFieldVal() throws NoSuchFieldException {
        Field field=MyObj.class.getDeclaredField("objField");
        long offset= U.objectFieldOffset(field);
        MyObj obj=new MyObj();

        int val= U.getInt(obj,offset);
        System.out.println("1.\t"+(val==10));

        MyObjChild child=new MyObjChild();
        int corVal1= U.getInt(child,offset);
        System.out.println("2.\t"+(corVal1==10));

        Field fieldChild=MyObj.class.getDeclaredField("objField");
        long offsetChild= U.objectFieldOffset(fieldChild);
        System.out.println("3.\t"+(offset==offsetChild));
        int corVal2= U.getInt(obj,offsetChild);
        System.out.println("4.\t"+(corVal2==10));

        short errVal1=U.getShort(obj,offset);
        System.out.println("5.\t"+(errVal1==10));

        int errVal2=U.getInt("abcd",offset);
        System.out.println("6.\t"+errVal2);

    }
输出结果为：

true
true
true
true
true
-223271518
第一个参数 o 和 offset 都是从 MyObj 获取的，所以返回 true。
第二个参数 o 是 MyObjChild 的实例，MyObjChild 是 MyObj 的子类，对象 o 是 MyObj 的兼容实例，所以返回 true。这从侧面说明在虚拟机中子类的实例的内存结构继承了父类的实例的内存结构。
第三个比较子类和父类中获取的字段偏移量是否相同，返回 true 说明是一样的，既然是一样的，第四个自然就返回 true。
这里重点说一下第五个，objField 是一个 int 类型，占四个字节，其值为 10，二进制为 00000000 00000000 00000000 00001010。Intel 处理器读取内存使用的是小端（Little-Endian）模式，在使用 Intel 处理器的机器的内存中多字节类型是按小端存储的，即低位在内存的低字节存储，高位在内存的高字节存储，所以 int 10 在内存中是（offset 0-3） 00001010 00000000 00000000 00000000。使用 getShort 会读取两个字节，即 00001010 00000000，获取的值仍为 10。

但是某些处理器是使用大端（Big-Endian），如 ARM 支持小端和大端，使用此处理器的机器的内存就会按大端存储多字节类型，与小端相反，此模式下低位在内存的高字节存储，高位在内存的低字节存储，所以 int 10 在内存中是（offset 0-3）00000000 00000000 00000000 00001010。在这种情况下，getShort 获取的值将会是 0。

不同的机器可能产生不一样的结果，基于此情况，如果字段是 int 类型，但需要一个 short 类型，也不应该调用 getShort，而应该调用 getInt，然后强制转换成 short。此外，如果调用 getLong，该方法返回的值一定不是 10。就像方法注释所说，调用该类型方法时，要保证方法的返回值和字段的值是同一种类型。

第五个测试获取非 MyObj 实例的偏移位置的值，这种情况下代码本身并不会报错，但获取到的值并非该字段的值（未定义的值）

2、offset 通过 staticFieldOffset 方法获取的类的某一字段的偏移量，o 是通过 staticFieldBase 方法获取的对象

    static void getClsFieldVal() throws NoSuchFieldException {
        Field field=MyObj.class.getDeclaredField("clsField");
        long offset= U.staticFieldOffset(field);
        Object obj=U.staticFieldBase(field);
        int val1=U.getInt(MyObj.class,offset);
        System.out.println("1.\t"+(val1==10));
        int val2=U.getInt(obj,offset);
        System.out.println("2.\t"+(val2==10));

    }
输出结果：

true
true
获取静态字段的值，有两个方法：staticFieldBase 获取字段所在的对象，静态字段附着于 Class 本身（java.lang.Class 的实例），该方法返回的其实就是该类本身，本例中是 MyObj.class。

3、如果 o 引用的是数组，则 offset 的值为 test.init.B+N*S，其中 N 是数组的合法的下标，test.init.B 是通过 arrayBaseOffset 方法从该数组类获取的，S 是通过 arrayIndexScale 方法从该数组类获取的

    static void getArrayVal(int index,int expectedVal) throws NoSuchFieldException {
        int base=U.arrayBaseOffset(int[].class);
        int scale=U.arrayIndexScale(int[].class);
        MyObj obj=new MyObj();
        Field field=MyObj.class.getDeclaredField("array");
        long offset= U.objectFieldOffset(field);
        int[] array= (int[]) U.getObject(obj,offset);
        int val1=U.getInt(array,(long)base+index*scale);
        System.out.println("1.\t"+(val1==expectedVal));
        int val2=U.getInt(obj.array,(long)base+index*scale);
        System.out.println("2.\t"+(val2==expectedVal));
    }
getArrayVal(2,30);
输出结果：

true
true
获取数组的值以及获取数组中某下标的值。获取数组某一下标的偏移量有一个计算公式 test.init.B+N*S，test.init.B 是数组元素在数组中的基准偏移量，S 是每个元素占用的字节数，N 是数组元素的下标。

有个要注意的地方，上面例子中方法内的数组的 offset 和 base 是两个完全不同的偏移量，offset 是数组 array 在对象 obj 中的偏移量，base 是数组元素在数组中的基准偏移量，这两个值没有任何联系，不能通过 offset 推导出 base。

getInt 的参数 o 可以是 null，在这种情况下，其和方法 getInt(long) 就是一样的了，offset 就不是表示相对的偏移地址了，而是表示内存中的绝对地址。操作系统中，一个进程是不能访问其他进程的内存的，所以传入 getInt 中的绝对地址必须是当前 JVM 管理的内存地址，否则进程会退出。

 

下一个方法，将值存储到 Java 变量中：

    public native void putInt(Object o, long offset, int x);
前两个参数会被解释成 Java 变量（字段或数组）的引用，
参数给定的值会被存储到该变量，变量的类型必须和方法参数的类型一致
参数 o 是变量关联的 Java 堆对象，可以为 null
参数 offset 代表该变量在该对象的位置，如果 o 是 null 则是内存的绝对地址
修改指定位置的内存，测试代码：

    static void setObjFieldVal(int val) throws NoSuchFieldException {
        Field field=MyObj.class.getDeclaredField("objField");
        long offset= U.objectFieldOffset(field);
        MyObj obj=new MyObj();
        U.putInt(obj,offset,val);
        int getVal= U.getInt(obj,offset);
        System.out.println(val==getVal);
        U.putLong(obj,offset,val);
        Field fieldArray=MyObj.class.getDeclaredField("array");
        long offsetArray= U.objectFieldOffset(fieldArray);
//        int[] array= (int[]) U.getObject(obj,offsetArray);
//        for(int i=0;i<array.length;i++){
//            System.out.println(array[i]);
//        }

    }
objField 是 int 类型，通过 putInt 修改 objField 的值可以正常修改，结果打印 true。然后使用 putLong 修改 objField 的值，这个修改操作本身也不会报错，但是 objField 并不是 long 类型，这样修改会导致其它程序错误，它不仅修改了 objField 的内存值，还修改了 objField 之后四个字节的内存值。

在这个例子中，objField 后面八个字节存储的是 array 字段所表示的数组对象的偏移位置，但是被修改了，如果后面的代码尝试访问 array 字段就会出错。修改其它字段（array、clsField）也是一样的，只要按之前的方法获取字段的偏移位置，使用与字段类型一致的 put 方法就可以。

下面的方法都是差不多的，只是针对不同的数据类型或者兼容 1.4 的字节码：

    public native void putObject(Object o, long offset, Object x);

    public native boolean getBoolean(Object o, long offset);

    public native void    putBoolean(Object o, long offset, boolean x);

    public native byte    getByte(Object o, long offset);

    public native void    putByte(Object o, long offset, byte x);

    public native short   getShort(Object o, long offset);
 
    public native void    putShort(Object o, long offset, short x);

    public native char    getChar(Object o, long offset);

    public native void    putChar(Object o, long offset, char x);

    public native long    getLong(Object o, long offset);

    public native void    putLong(Object o, long offset, long x);

    public native float   getFloat(Object o, long offset);

    public native void    putFloat(Object o, long offset, float x);

    public native double  getDouble(Object o, long offset);

    public native void    putDouble(Object o, long offset, double x);

    @Deprecated
    public int getInt(Object o, int offset) {
        return getInt(o, (long)offset);
    }

   
    @Deprecated
    public void putInt(Object o, int offset, int x) {
        putInt(o, (long)offset, x);
    }

   
    @Deprecated
    public Object getObject(Object o, int offset) {
        return getObject(o, (long)offset);
    }

   
    @Deprecated
    public void putObject(Object o, int offset, Object x) {
        putObject(o, (long)offset, x);
    }

    
    @Deprecated
    public boolean getBoolean(Object o, int offset) {
        return getBoolean(o, (long)offset);
    }

   
    @Deprecated
    public void putBoolean(Object o, int offset, boolean x) {
        putBoolean(o, (long)offset, x);
    }

    
    @Deprecated
    public byte getByte(Object o, int offset) {
        return getByte(o, (long)offset);
    }

    
    @Deprecated
    public void putByte(Object o, int offset, byte x) {
        putByte(o, (long)offset, x);
    }

    
    @Deprecated
    public short getShort(Object o, int offset) {
        return getShort(o, (long)offset);
    }

   
    @Deprecated
    public void putShort(Object o, int offset, short x) {
        putShort(o, (long)offset, x);
    }

    
    @Deprecated
    public char getChar(Object o, int offset) {
        return getChar(o, (long)offset);
    }

    
    @Deprecated
    public void putChar(Object o, int offset, char x) {
        putChar(o, (long)offset, x);
    }

    
    @Deprecated
    public long getLong(Object o, int offset) {
        return getLong(o, (long)offset);
    }

    
    @Deprecated
    public void putLong(Object o, int offset, long x) {
        putLong(o, (long)offset, x);
    }

    
    @Deprecated
    public float getFloat(Object o, int offset) {
        return getFloat(o, (long)offset);
    }

    
    @Deprecated
    public void putFloat(Object o, int offset, float x) {
        putFloat(o, (long)offset, x);
    }

    
    @Deprecated
    public double getDouble(Object o, int offset) {
        return getDouble(o, (long)offset);
    }

    
    @Deprecated
    public void putDouble(Object o, int offset, double x) {
        putDouble(o, (long)offset, x);
    }
下面的方法和上面的也类似，只是这些方法只有一个参数，即内存绝对地址，这些方法不需要 Java 对象地址作为基准地址，所以它们可以作用于本地方法区：


    //获取内存地址的值，如果地址是 0 或者不是指向通过 allocateMemory 方法获取的内存块，则结果是未知的
    
    public native byte    getByte(long address);


    //将一个值写入内存，如果地址是 0 或者不是指向通过 allocateMemory 方法获取的内存块，则结果是未知的

    public native void    putByte(long address, byte x);

    public native short   getShort(long address);

    public native void    putShort(long address, short x);

    public native char    getChar(long address);

    public native void    putChar(long address, char x);

    public native int     getInt(long address);

    public native void    putInt(long address, int x);

    public native long    getLong(long address);
 
    public native void    putLong(long address, long x);

    public native float   getFloat(long address);
 
    public native void    putFloat(long address, float x);
  
    public native double  getDouble(long address);
  
    public native void    putDouble(long address, double x);
这里提到一个方法 allocateMemory，它是用于分配本地内存的。看看和本地内存有关的几个方法：

    ///包装malloc，realloc，free

    /**
     * 分配指定大小的一块本地内存。分配的这块内存不会初始化，它们的内容通常是没用的数据
     * 返回的本地指针不会是 0，并且该内存块是连续的。调用 freeMemory 方法可以释放此内存，调用
     * reallocateMemory 方法可以重新分配
     */
    public native long allocateMemory(long bytes);

    /**
     * 重新分配一块指定大小的本地内存，超出老内存块的字节不会被初始化，它们的内容通常是没用的数据
     * 当且仅当请求的大小为 0 时，该方法返回的本地指针会是 0。
     * 该内存块是连续的。调用 freeMemory 方法可以释放此内存，调用 reallocateMemory 方法可以重新分配
     * 参数 address 可以是 null，这种情况下会分配新内存(和 allocateMemory 一样)
     */
    public native long reallocateMemory(long address, long bytes);


    /** 
     * 将给定的内存块的所有字节设置成固定的值(通常是 0)
     * 该方法通过两个参数确定内存块的基准地址，就像在 getInt(Object,long) 中讨论的，它提供了 double-register 地址模型
     * 如果引用的对象是 null, 则 offset 会被当成绝对基准地址
     * 该写入操作是按单元写入的，单元的字节大小由地址和长度参数决定，每个单元的写入是原子性的。如果地址和长度都是 8 的倍数，则一个单元为 long
     * 型(一个单元 8 个字节)；如果地址和长度都是 4 的倍数，则一个单元为 int 型(一个单元 4 个字节)；
     * 如果地址和长度都是 2 的倍数，则一个单元为 short 型(一个单元 2 个字节)；
     */

    public native void setMemory(Object o, long offset, long bytes, byte value);


    //将给定的内存块的所有字节设置成固定的值(通常是 0)
    //就像在 getInt(Object,long) 中讨论的，该方法提供 single-register 地址模型

    public void setMemory(long address, long bytes, byte value) {
        setMemory(null, address, bytes, value);
    }


    //复制指定内存块的字节到另一内存块
    //该方法的两个基准地址分别由两个参数决定

    public native void copyMemory(Object srcBase, long srcOffset,
                                  Object destBase, long destOffset,
                                  long bytes);

    //复制指定内存块的字节到另一内存块，但使用 single-register 地址模型
     
    public void copyMemory(long srcAddress, long destAddress, long bytes) {
        copyMemory(null, srcAddress, null, destAddress, bytes);
    }

    //释放通过 allocateMemory 或者 reallocateMemory 获取的内存，如果参数 address 是 null，则不做任何处理

    public native void freeMemory(long address);
allocateMemory、reallocateMemory、freeMemory 与 setMemory 分别是对 C 函数 malloc、realloc、free 和 memset 的封装，这样该类就提供了动态获取/释放本地方法区内存的功能。

malloc 用于分配一个全新的未使用的连续内存，但该内存不会初始化，即不会被清零；
realloc 用于内存的缩容或扩容，有两个参数，从 malloc 返回的地址和要调整的大小，该函数和 malloc 一样，不会初始化，它能保留之前放到内存里的值，很适合用于扩容；
free 用于释放内存，该方法只有一个地址参数，那它如何知道要释放多少个字节呢？其实在 malloc 分配内存的时候会多分配 4 个字节用于存放该块的长度，比如 malloc(10) 其实会花费 14 个字节。理论上讲能分配的最大内存是 4G（2^32-1）。在 hotspot 虚拟机的设计中，数组对象也有 4 个字节用于存放数组长度，那么在 hotspot 中，数组的最大长度就是 2^32-1，这样 free 函数只要读取前 4 个字节就知道要释放多少内存了（10+4）；
memset  一般用于初始化内存，可以设置初始化内存的值，一般初始值会设置成 0，即清零操作。
来一个简单的例子，申请内存-写入内存-读取内存-释放内存：

long address=U.allocateMemory(10);
U.setMemory(address,10,(byte)1);
/**
 * 1的二进制码为00000001，int为四个字节，U.getInt将读取四个字节，
 * 读取的字节为00000001 00000001 00000001 00000001
 */
int i=0b00000001000000010000000100000001;
System.out.println(i==U.getInt(address));
U.freeMemory(address);
 

接下来看看获取类变量相关信息的几个方法：

    /// random queries
    /// 随机搜索，对象是放在一块连续的内存空间中，所以是支持随机搜索的
    /**
     * staticFieldOffset，objectFieldOffset，arrayBaseOffset 方法的返回值不会是该常量(-1)
     */
    public static final int INVALID_FIELD_OFFSET   = -1;

    /**
     * 返回字段的偏移量，32 字节
     * 从 1.4.1 开始，对于静态字段，请使用 staticFieldOffset 方法，非静态字段使用 objectFieldOffset 方法获取
     */
    @Deprecated
    public int fieldOffset(Field f) {
        if (Modifier.isStatic(f.getModifiers()))
            return (int) staticFieldOffset(f);
        else
            return (int) objectFieldOffset(f);
    }

    /**
     * 返回用于访问静态字段的基准地址
     * 从 1.4.1 开始，要获取访问指定字段的基准地址，请使用 staticFieldBase(Field)
     * 该方法仅能作用于把所有静态字段放在一起的 JVM 实现
     */
    @Deprecated
    public Object staticFieldBase(Class<?> c) {
        Field[] fields = c.getDeclaredFields();
        for (int i = 0; i < fields.length; i++) {
            if (Modifier.isStatic(fields[i].getModifiers())) {
                return staticFieldBase(fields[i]);
            }
        }
        return null;
    }

    /**
     * 返回给定的字段在该类的偏移地址
     *
     * 对于任何给定的字段，该方法总是返回相同的值；同一个类的不同字段总是返回不同的值
     *
     * 从 1.4.1 开始，字段的偏移以 long 表示，虽然 Sun 的 JVM 只使用了 32 位，但是那些将静态字段存储到绝对地址的 JVM 实现
     * 需要使用 long 类型的偏移量，通过 getXX(null,long) 获取字段值，为了保持代码迁移到 64 位平台上 JVM 的优良性，
     * 必须保持静态字段偏移量的所有比特位
     */
    public native long staticFieldOffset(Field f);

    /**
     * 很难想象 JVM 需要使用这么多比特位来编码非数组对象的偏移量，它们只需要很少的比特位就可以了(有谁看过有100个成员变量的类么？
     * 一个字节能表示 256 个成员变量)，
     * 为了和该类的其他方法保持一致，所以该方法也返回 long 类型
     *
     */
    public native long objectFieldOffset(Field f);

    /**
     * 获取指定静态字段的位置，和 staticFieldOffset 一起使用
     * 获取该静态字段所在的"对象"，这个"对象"可通过类似 getInt(Object,long) 的方法访问
     * 该"对象"可能是 null，并且引用的可能是对象的"cookie"(此处cookie具体含义未知，没有找到相关资料)，不保证是真正的对象，该"对象"只能当作此类中 put 和 get 方法的参数，
     * 其他情况下不应该使用它
     */
    public native Object staticFieldBase(Field f);

    /**
     * 检查给定的类是否需要初始化，它通常和 staticFieldBase 方法一起使用
     * 只有当 ensureClassInitialized 方法不产生任何影响时才会返回 false
     */
    public native boolean shouldBeInitialized(Class<?> c);

    /**
     * 确保给定的类已被初始化，它通常和 staticFieldBase 方法一起使用
     */
    public native void ensureClassInitialized(Class<?> c);

    /**
     * 返回给定数组类第一个元素在内存中的偏移量，如果 arrayIndexScale 方法返回非0值，要获得访问数组元素的新的偏移量，
     * 需要使用 scale
     */
    public native int arrayBaseOffset(Class<?> arrayClass);

    /** The value of {@code arrayBaseOffset(boolean[].class)} */
    public static final int ARRAY_BOOLEAN_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(boolean[].class);

    /** The value of {@code arrayBaseOffset(byte[].class)} */
    public static final int ARRAY_BYTE_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(byte[].class);

    /** The value of {@code arrayBaseOffset(short[].class)} */
    public static final int ARRAY_SHORT_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(short[].class);

    /** The value of {@code arrayBaseOffset(char[].class)} */
    public static final int ARRAY_CHAR_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(char[].class);

    /** The value of {@code arrayBaseOffset(int[].class)} */
    public static final int ARRAY_INT_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(int[].class);

    /** The value of {@code arrayBaseOffset(long[].class)} */
    public static final int ARRAY_LONG_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(long[].class);

    /** The value of {@code arrayBaseOffset(float[].class)} */
    public static final int ARRAY_FLOAT_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(float[].class);

    /** The value of {@code arrayBaseOffset(double[].class)} */
    public static final int ARRAY_DOUBLE_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(double[].class);

    /** The value of {@code arrayBaseOffset(Object[].class)} */
    public static final int ARRAY_OBJECT_BASE_OFFSET
            = theUnsafe.arrayBaseOffset(Object[].class);

    /**
     * 返回给定数组类的每个元素在内存中的 scale(所占用的字节)。然而对于"narrow"类型的数组，类似 getByte(Object, int) 的访问方法
     * 一般不会获得正确的结果，所以这些类返回的 scale 会是 0
     * (本人水平有限，此处narrow类型不知道具体含义，不了解什么时候此方法会返回0)
     */
    public native int arrayIndexScale(Class<?> arrayClass);

    /** The value of {@code arrayIndexScale(boolean[].class)} */
    public static final int ARRAY_BOOLEAN_INDEX_SCALE
            = theUnsafe.arrayIndexScale(boolean[].class);

    /** The value of {@code arrayIndexScale(byte[].class)} */
    public static final int ARRAY_BYTE_INDEX_SCALE
            = theUnsafe.arrayIndexScale(byte[].class);

    /** The value of {@code arrayIndexScale(short[].class)} */
    public static final int ARRAY_SHORT_INDEX_SCALE
            = theUnsafe.arrayIndexScale(short[].class);

    /** The value of {@code arrayIndexScale(char[].class)} */
    public static final int ARRAY_CHAR_INDEX_SCALE
            = theUnsafe.arrayIndexScale(char[].class);

    /** The value of {@code arrayIndexScale(int[].class)} */
    public static final int ARRAY_INT_INDEX_SCALE
            = theUnsafe.arrayIndexScale(int[].class);

    /** The value of {@code arrayIndexScale(long[].class)} */
    public static final int ARRAY_LONG_INDEX_SCALE
            = theUnsafe.arrayIndexScale(long[].class);

    /** The value of {@code arrayIndexScale(float[].class)} */
    public static final int ARRAY_FLOAT_INDEX_SCALE
            = theUnsafe.arrayIndexScale(float[].class);

    /** The value of {@code arrayIndexScale(double[].class)} */
    public static final int ARRAY_DOUBLE_INDEX_SCALE
            = theUnsafe.arrayIndexScale(double[].class);

    /** The value of {@code arrayIndexScale(Object[].class)} */
    public static final int ARRAY_OBJECT_INDEX_SCALE
            = theUnsafe.arrayIndexScale(Object[].class);
上面的一些方法之前已经提到过，注释也说的比较明白, 说一下 shouldBeInitialized 和 ensureClassInitialized，shouldBeInitialized 判断类是否已初始化，ensureClassInitialized 执行初始化。有个概念需要了解，虚拟机加载类包括加载和链接阶段，加载阶段只是把类加载进内存，链接阶段会验证加载的代码的合法性，并初始化静态字段和静态块；shouldBeInitialized 就是检查链接阶段有没有执行。

    static void clsInitialized() throws NoSuchFieldException {
        System.out.println(U.shouldBeInitialized(MyObj.class));
        System.out.println(U.shouldBeInitialized(MyObjChild.class));
        U.ensureClassInitialized(MyObjChild.class);
        System.out.println(U.shouldBeInitialized(MyObjChild.class));
    }
public class MyObjChild extends MyObj {
    static int f1=1;
    final static int f2=1;
    static {
        f1=2;
        System.out.println("MyObjChild init");
    }
}
输出：

false 
true 
MyObjChild init 
false
第一行输出 false 是因为我这个代码（包括 main 方法）是在 MyObj 类里写的，执行 main 的时候，MyObj 已经加载并初始化了。调用 U.shouldBeInitialized(MyObjChild.class) 只会加载 MyObjChild.class，但不会初始化，执行 ensureClassInitialized 才会初始化。

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
输出：

1.true
2.true MyObjChild init
3.true 
f1 是 static int，f2 是 final static int，因为 f2 是 final，它的值在编译期就决定了，存放在类的常量表里，所以即使还没有初始化它的值就是 1。

 

    /**
     * 获取本地指针所占用的字节大小，值为 4 或者 8。其他基本类型的大小由其内容决定
     */
    public native int addressSize();

    /** The value of {@code addressSize()} */
    public static final int ADDRESS_SIZE = theUnsafe.addressSize();

    /**
     * 本地内存页大小，值为 2 的 N 次方
     */
    public native int pageSize();
addressSize 返回指针的大小，32 位虚拟机返回 4，64 位虚拟机默认返回 8，开启指针压缩功能（-XX:-UseCompressedOops）则返回 4。基本类型不是用指针表示的，它是直接存储的值。一般情况下，我们会说在 Java 中，基本类型是值传递，对象是引用传递。Java 官方的表述是在任何情况下 Java 都是值传递。基本类型是传递值本身，对象类型是传递指针的值。

 

    /// random trusted operations from JNI:
    /// JNI信任的操作
    /**
     * 告诉虚拟机定义一个类，加载类不做安全检查，默认情况下，参数类加载器(ClassLoader)和保护域(ProtectionDomain)来自调用者类
     */
    public native Class<?> defineClass(String name, byte[] b, int off, int len,
                                       ClassLoader loader,
                                       ProtectionDomain protectionDomain);

    /**
     * 定义一个匿名类，这里说的和我们代码里写的匿名内部类不是一个东西。
     * (可以参考知乎上的一个问答 https://www.zhihu.com/question/51132462)
     */
    public native Class<?> defineAnonymousClass(Class<?> hostClass, byte[] data, Object[] cpPatches);


    /** 
     * 分配实例的内存空间，但不会执行构造函数。如果没有执行初始化，则会执行初始化
     */
    public native Object allocateInstance(Class<?> cls)
            throws InstantiationException;

    /** Lock the object.  It must get unlocked via {@link #monitorExit}.
     *
     * 获取对象内置锁(即 synchronized 关键字获取的锁)，必须通过 monitorExit 方法释放锁
     * (synchronized 代码块在编译后会产生两个指令:monitorenter,monitorexit)
     */
    public native void monitorEnter(Object o);

    /**
     * Unlock the object.  It must have been locked via {@link
     * #monitorEnter}.
     * 释放锁
     */
    public native void monitorExit(Object o);

    /**
     * 尝试获取对象内置锁，通过返回 true 和 false 表示是否成功获取锁
     */
    public native boolean tryMonitorEnter(Object o);

    /** Throw the exception without telling the verifier.
     * 不通知验证器(verifier)直接抛出异常(此处 verifier 具体含义未知，没有找到相关资料)
     */
    public native void throwException(Throwable ee);
    allocateInstance 方法的测试
    
    public class MyObjChild extends MyObj {
        static int f1=1;
        int f2=1;
        static {
            f1=2;
            System.out.println("MyObjChild init");
        }
        public MyObjChild(){
            f2=2;
            System.out.println("run construct");
        }
    }
       static void clsInitialized3() throws InstantiationException {
            MyObjChild myObj= (MyObjChild) U.allocateInstance(MyObjChild.class);
            System.out.println("1.\t"+(MyObjChild.f1==2));
            System.out.println("1.\t"+(myObj.f2==0));
        }
输出：

MyObjChild init
1.true
2.true
可以看到分配对象的时候只执行了类的初始化代码，没有执行构造函数。

 

来看看最重要的 CAS 方法

    /**
     * Atomically update Java variable to <tt>x</tt> if it is currently
     * holding <tt>expected</tt>.
     * @return <tt>true</tt> if successful
     *
     * 如果变量的值为预期值，则更新变量的值，该操作为原子操作
     * 如果修改成功则返回true
     */
    public final native boolean compareAndSwapObject(Object o, long offset,
                                                     Object expected,
                                                     Object x);

    /**
     * Atomically update Java variable to <tt>x</tt> if it is currently
     * holding <tt>expected</tt>.
     * @return <tt>true</tt> if successful
     */
    public final native boolean compareAndSwapInt(Object o, long offset,
                                                  int expected,
                                                  int x);

    /**
     * Atomically update Java variable to <tt>x</tt> if it is currently
     * holding <tt>expected</tt>.
     * @return <tt>true</tt> if successful
     */
    public final native boolean compareAndSwapLong(Object o, long offset,
                                                   long expected,
                                                   long x);
这几个方法应该是最常用的方法了，用于实现原子性的 CAS 操作，这些操作可以避免加锁，一般情况下，性能会更好， java.util.concurrent 包下很多类就是用的这些 CAS 操作而没有用锁。

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
输出：

1.false
2.true
 

    /**
     * 获取给定变量的引用值，该操作有 volatile 加载语意，其他方面和 getObject(Object, long) 一样
     */
    public native Object getObjectVolatile(Object o, long offset);

    /**
     * 将引用值写入给定的变量，该操作有 volatile 加载语意，其他方面和 putObject(Object, long, Object) 一样
     */
    public native void    putObjectVolatile(Object o, long offset, Object x);

    /** Volatile version of {@link #getInt(Object, long)}  */
    public native int     getIntVolatile(Object o, long offset);

    /** Volatile version of {@link #putInt(Object, long, int)}  */
    public native void    putIntVolatile(Object o, long offset, int x);

    /** Volatile version of {@link #getBoolean(Object, long)}  */
    public native boolean getBooleanVolatile(Object o, long offset);

    /** Volatile version of {@link #putBoolean(Object, long, boolean)}  */
    public native void    putBooleanVolatile(Object o, long offset, boolean x);

    /** Volatile version of {@link #getByte(Object, long)}  */
    public native byte    getByteVolatile(Object o, long offset);

    /** Volatile version of {@link #putByte(Object, long, byte)}  */
    public native void    putByteVolatile(Object o, long offset, byte x);

    /** Volatile version of {@link #getShort(Object, long)}  */
    public native short   getShortVolatile(Object o, long offset);

    /** Volatile version of {@link #putShort(Object, long, short)}  */
    public native void    putShortVolatile(Object o, long offset, short x);

    /** Volatile version of {@link #getChar(Object, long)}  */
    public native char    getCharVolatile(Object o, long offset);

    /** Volatile version of {@link #putChar(Object, long, char)}  */
    public native void    putCharVolatile(Object o, long offset, char x);

    /** Volatile version of {@link #getLong(Object, long)}  */
    public native long    getLongVolatile(Object o, long offset);

    /** Volatile version of {@link #putLong(Object, long, long)}  */
    public native void    putLongVolatile(Object o, long offset, long x);

    /** Volatile version of {@link #getFloat(Object, long)}  */
    public native float   getFloatVolatile(Object o, long offset);

    /** Volatile version of {@link #putFloat(Object, long, float)}  */
    public native void    putFloatVolatile(Object o, long offset, float x);

    /** Volatile version of {@link #getDouble(Object, long)}  */
    public native double  getDoubleVolatile(Object o, long offset);

    /** Volatile version of {@link #putDouble(Object, long, double)}  */
    public native void    putDoubleVolatile(Object o, long offset, double x);
这是具有 volatile 语意的 get 和 put方法。volatile 语意为保证不同线程之间的可见行，即一个线程修改一个变量之后，保证另一线程能观测到此修改。这些方法可以使非 volatile 变量具有 volatile 语意。

 

    /**
     * putObjectVolatile(Object, long, Object)的另一个版本(有序的/延迟的)，它不保证其他线程能立即看到修改,
     * 该方法通常只对底层为 volatile 的变量(或者 volatile 类型的数组元素)有帮助
     */
    public native void    putOrderedObject(Object o, long offset, Object x);

    /** Ordered/Lazy version of {@link #putIntVolatile(Object, long, int)}  */
    public native void    putOrderedInt(Object o, long offset, int x);

    /** Ordered/Lazy version of {@link #putLongVolatile(Object, long, long)} */
    public native void    putOrderedLong(Object o, long offset, long x);
有三类很相近的方法：putXx、putXxVolatile 与 putOrderedXx：

putXx 只是写本线程缓存，不会将其它线程缓存置为失效，所以不能保证其它线程一定看到此次修改；
putXxVolatile 相反，它可以保证其它线程一定看到此次修改；
putOrderedXx 也不保证其它线程一定看到此次修改，但和 putXx 又有区别，它的注释上有两个关键字：顺序性（Ordered）和延迟性（lazy），顺序性是指不会发生重排序，延迟性是指其它线程不会立即看到此次修改，只有当调用 putXxVolatile 使才能看到。
 

    /**
     * 释放当前阻塞的线程。如果当前线程没有阻塞，则下一次调用 park 不会阻塞。这个操作是"非安全"的
     * 是因为调用者必须通过某种方式保证该线程没有被销毁
     *
     */
    public native void unpark(Object thread);

    /**
     * 阻塞当前线程，当发生如下情况时返回：
     * 1、调用 unpark 方法
     * 2、线程被中断
     * 3、时间过期
     * 4、spuriously
     * 该操作放在 Unsafe 类里没有其它意义，它可以放在其它的任何地方
     */
    public native void park(boolean isAbsolute, long time);
阻塞和释放当前线程，java.util.concurrent 中的锁就是通过这两个方法实现线程阻塞和释放的。

 

    /**
     *获取一段时间内，运行的任务队列分配到可用处理器的平均数(平常说的 CPU 使用率)
     *
     */
    public native int getLoadAverage(double[] loadavg, int nelems);
统计 CPU 负载。

 

    // The following contain CAS-based Java implementations used on
    // platforms not supporting native instructions
    //下面的方法包含基于 CAS 的 Java 实现，用于不支持本地指令的平台
    /**
     * 在给定的字段或数组元素的当前值原子性的增加给定的值
     * @param o 字段/元素所在的对象/数组
     * @param offset 字段/元素的偏移
     * @param delta 需要增加的值
     * @return 原值
     * @since 1.8
     */
    public final int getAndAddInt(Object o, long offset, int delta) {
        int v;
        do {
            v = getIntVolatile(o, offset);
        } while (!compareAndSwapInt(o, offset, v, v + delta));
        return v;
    }

    public final long getAndAddLong(Object o, long offset, long delta) {
        long v;
        do {
            v = getLongVolatile(o, offset);
        } while (!compareAndSwapLong(o, offset, v, v + delta));
        return v;
    }

    /**
     * 将给定的字段或数组元素的当前值原子性的替换给定的值
     * @param o 字段/元素所在的对象/数组
     * @param offset field/element offset
     * @param newValue 新值
     * @return 原值
     * @since 1.8
     */
    public final int getAndSetInt(Object o, long offset, int newValue) {
        int v;
        do {
            v = getIntVolatile(o, offset);
        } while (!compareAndSwapInt(o, offset, v, newValue));
        return v;
    }

    public final long getAndSetLong(Object o, long offset, long newValue) {
        long v;
        do {
            v = getLongVolatile(o, offset);
        } while (!compareAndSwapLong(o, offset, v, newValue));
        return v;
    }

    public final Object getAndSetObject(Object o, long offset, Object newValue) {
        Object v;
        do {
            v = getObjectVolatile(o, offset);
        } while (!compareAndSwapObject(o, offset, v, newValue));
        return v;
    }
基于 CAS 的一些原子操作实现，也是比较常用的方法。

 

    //确保该栏杆前的读操作不会和栏杆后的读写操作发生重排序

    public native void loadFence();

    //确保该栏杆前的写操作不会和栏杆后的读写操作发生重排序

    public native void storeFence();

    //确保该栏杆前的读写操作不会和栏杆后的读写操作发生重排序

    public native void fullFence();

    //抛出非法访问错误，仅用于VM内部

    private static void throwIllegalAccessError() {
        throw new IllegalAccessError();
    }
这是实现内存屏障的几个方法，类似于 volatile 的语意，保证内存可见性和禁止重排序。这几个方法涉及到 JMM（Java 内存模型），有兴趣的可参考Java 内存模型 Cookbook 翻译 。

作者介绍
相奕，互联网开发者，多年互联网开发经验，关注底层技术与微服务周边技术。

本文系作者投稿文章，欢迎投稿。投稿要求见：

https://my.oschina.net/editorial-story/blog/1814725


很好的文章, 我有些补充:
1. JDK11包含了Unsafe源码,内部实现转移到jdk.internal.misc这个包里,仍然属于半开放状态,提供的方法可能随时修改.
2. JDK11的Unsafe大部分方法都带@HotSpotIntrinsicCandidate注解,说明会被直接映射成机器指令,而不是方法调用. 还有些方法用了@ForceInline注解,能保证被内联优化.
3. "Field fieldChild=MyObj.class.getDeclaredField("objField");"里的MyObj.class应该是MyObjChild.class.
4. "第五个测试获取非 MyObj 实例的偏移位置的值"应该是第六个.
5. 文中提到的几个@Deprecated方法在JDK11里不是删除了就是原型改变了(主要是offset变成了long类型)
6. "narrow类型"估计是指数组中的每个元素占用不到1个字节的情况, 目前所见的JDK实现没有这种情况, 估计有的JDK把boolean数组中的每个元素只用1 bit来实现的.
7. "Java 官方的表述是在任何情况下 Java 都是值传递。基本类型是传递值本身，对象类型是传递指针的值。"我完全赞同这种说法,非常清晰准确,我觉得严格的"引用"含义应该像C++中的"&"和C#中的"ref/out"这种变量本身的别名(就是说Java对象的引用应该是指针的指针).
8. "不通知验证器(verifier)直接抛出异常"应该是指可以绕过throws声明.
