package xyz.hoper.test.lock;

/**
 * @Description TODO
 * @Date 2022/2/21 17:55
 * @Created by lbyi
 */
public class Lock {

   synchronized void method1(boolean operator){
     if(operator) Num.count++;
       else Num.count--;
     System.out.print(Num.count);
  }
  void method2(boolean operator){
    if(operator) Num.count++;
    else Num.count--;
    System.out.print(Num.count);
  }
}

