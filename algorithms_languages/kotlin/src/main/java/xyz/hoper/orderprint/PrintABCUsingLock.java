package xyz.hoper.orderprint;

import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;

/**
 * @author ：lbyi
 * @date ：Created in 2019/3/11 13:53
 * @description：多线程打印
 * @modified By：
 */
public class PrintABCUsingLock {
    private final int times;
    private int state;
    private final Lock lock = new ReentrantLock();

    public PrintABCUsingLock(int times) {
        this.times = times;
    }

    public static void main(String[] args) {
        PrintABCUsingLock printABC = new PrintABCUsingLock(1000);
        //非静态方法引用 x::toString == ()->x.toString()
        new Thread(printABC::printA).start();
        new Thread(printABC::printB).start();
        new Thread(printABC::printC).start();
        //new Thread(()->printABC.printA()).start();
    }

    public void printA() {
        print("A", 0);
    }

    public void printB() {
        print("B", 1);
    }

    public void printC() {
        print("C", 2);
    }

    private void print(String name, int targetState) {
        for (int i = 0; i < times; ) {
            lock.lock();
            if (state % 3 == targetState) {
                state++;
                i++;
                System.out.println(name + i);
            }
            System.out.println(name + "空循环");
            lock.unlock();
        }
    }
}
