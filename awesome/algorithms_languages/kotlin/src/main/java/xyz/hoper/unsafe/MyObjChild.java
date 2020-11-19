package xyz.hoper.unsafe;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 9:32
 * @description：
 * @modified By：
 */
public class MyObjChild extends MyObj {
    static int f1=1;
    int f2=1;
    final static int f3 =1;
    static {
        f1=2;
        System.out.println("MyObjChild init");
    }

    public MyObjChild(){
        f2=2;
        System.out.println("run construct");
    }
}
