package xyz.hoper.test.lock;

import java.util.concurrent.locks.ReentrantLock;

import java.lang.Thread;

public class Num {
  static int count;

  public static void main(String[] args) throws InterruptedException {
    var l1 = new Lock();
    var l2 = new Lock();
    for (int i = 0; i <10;i++) new Thread(() -> {
       l1.method1(true);
    }).start();
    for (int i = 0; i <10;i++) new Thread(() -> {
       l2.method1(false);
    }).start();
    Thread.sleep(1000);
    final java.util.concurrent.locks.Lock lock = new ReentrantLock();
    for (int i = 0; i <10;i++) new Thread(() -> {
      lock.lock();
      l1.method2(true);
      lock.unlock();
    }).start();
    for (int i = 0; i <10;i++) new Thread(() -> {
      lock.lock();
      l1.method2(false);
      lock.unlock();
    }).start();
  }
}
