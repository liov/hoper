由于c，go，java的内存模型不一致，所以相互调用会有额外开销，所以为了炫技强行上cgo和jni开有时候得不偿失，为了知道cgo和jni的应用场景，特意做了测试

事先想到，c比go绝对的快，比java快，那么调用c的应用场景自然是cpu密集型，废话不多说，上斐波那契

cgo的代码

```golang
    /*
unsigned long long int fibonacci(unsigned int n);


unsigned long long int fibonacci(unsigned int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

 */
import "C"
import (
   "fmt"
   "time"
)
//可能这就是cgo的应用场景吧，但是真的存在吗，n为34时，时间一致，n为42时，时间差了一倍
//go向c传参有额外开销
func main()  {
   useGo()
   useCGo()
}

func useGo()  {
   now:=time.Now()
   fmt.Println(now)
   fmt.Println(C.fibonacci(42))
   fmt.Println(time.Now().Sub(now))
}

func useCGo()  {
   now:=time.Now()
   fmt.Println(now)
   fmt.Println(goFib(42))
   fmt.Println(time.Now().Sub(now))
}

func goFib(n uint) uint64 {
   if n <2 {return 1}
   return goFib(n - 1) + goFib(n - 2)
}
```
经过测试，go不向c传参时，n为34是个分界岭，此前go运行时间小于cgo，此后cgo小于go，甚至达到一倍的差距

go传参时，n为42时，cgo耗时1s，go耗时1.8s，差距还是很明显的

由此得出，cgo还是有应用场景的

 

接下来看jni，上代码
```java
public class JNI {

    //链动态库
    static {
        System.loadLibrary("hello");
    }

    //方法定义
    public native void testHelloVoid();

    public native String testHello();

    public static native long fib(int n);

    public static void main(String[] args){
        //执行
        jnifib();
        javafib();
    }

    private static void jnifib(){
        long starTime=System.currentTimeMillis();
        System.out.println(starTime);
        System.out.println(fib(43));
        System.out.println(System.currentTimeMillis()-starTime);
    }

    private static void javafib(){
        long starTime=System.currentTimeMillis();
        System.out.println(starTime);
        System.out.println(jfib(43));
        System.out.println(System.currentTimeMillis()-starTime);
    }

    private static long jfib(int n){
        if(n<2) return 1;
        return jfib(n-1)+jfib(n-2);
    }

}
```

编译步骤就不详述了，网上很多

jni的测试结果是，java的运行时间总小于jni，且随着n的增大，时间差也更大，并没有发生所谓分界

猜想原因，首先java优化足够好了，跟C的差距不大，其实，我估计是调用开销过大，跟go的调用逻辑应该不一样，只调用一次C函数并不是C执行然后读内存读值，具体逻辑得google了

