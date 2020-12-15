package xyz.hoper.test.jni;

/**
 * @author ：lbyi
 * @date ：Created in 2019/5/13
 * @description：jni
 */

//javac -h . JNI.java or javac JNI.java -h JniH
//java jni的开销只会越来越大,看来java优化足够好了，不需要c了
//上一句是错的，原因是编译出来的C库没开优化
public class JNI {

    //链动态库
    static {
        System.loadLibrary("hello");
    }

    //方法定义
    public native void testHelloVoid();

    public native String testHello();

    public static native long fib(int n);

    public static native void jnifor();

    public static void main(String[] args){
        //执行
        jnifib();
        javafib();
        jjnifor();
        jfor();
    }

    private static void jnifib(){
        long starTime=System.currentTimeMillis();
        System.out.println("jni:"+starTime);
        System.out.println(fib(42));
        System.out.println(System.currentTimeMillis()-starTime);
    }

    private static void javafib(){
        long starTime=System.currentTimeMillis();
        System.out.println("java:"+starTime);
        System.out.println(jfib(42));
        System.out.println(System.currentTimeMillis()-starTime);
    }

    private static long jfib(int n){
        if(n<2) return 1;
        return jfib(n-1)+jfib(n-2);
    }

    private static void jjnifor(){
        long starTime=System.currentTimeMillis();
        System.out.println("jnifor:"+starTime);
        jnifor();
        System.out.println(System.currentTimeMillis()-starTime);
    }

    private  static void jfor(){
        long starTime=System.currentTimeMillis();
        System.out.println("jfor:"+starTime);
        long j=0;
        for(long i=0;i<100000000;i++){
           j++;
        }
        System.out.println(System.currentTimeMillis()-starTime);
    }

}
