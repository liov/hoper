package xyz.hoper.unsafe;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/26 9:34
 * @description：
 * @modified By：
 */
public class UnsafeTest {

    public static void main(String[] args){

        try {
            MyObj.getObjFieldVal();
            MyObj.getClsFieldVal();
            MyObj.getArrayVal(2,30);
            MyObj.setObjFieldVal(15);
        } catch (NoSuchFieldException e) {
            e.printStackTrace();
        }

        MyObj.memory();
    }

}
