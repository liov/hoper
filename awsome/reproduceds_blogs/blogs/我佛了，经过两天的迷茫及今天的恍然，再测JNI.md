上篇提到，JNI的反常情况，调用JNI反而更慢

两天来搜各种文献均没有得到有用资料，一度迷茫

今天忽然想到，在C中看C的用时

结果出乎意料

```c
int begintime,endtime;
begintime = clock();
long long int i = fibonacci(n);
endtime = clock();
printf("Running Time:%f ms\n", (double)(endtime-begintime));
```
测试结果显示，C中调用时间和JNI总耗时接近，且是相当接近，但是这就表明一个问题，C比java慢，？？？

不行，我要做纯C测试

测试结果让人大跌眼镜，C比java慢，至少是fib上，代码也很简单
```c
long long int fibonacci(int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}
```
C水平不佳，想不到优化点

但java性能也不至于恐怖如斯吧，我得到linux上再测一遍

为什么同一个版本的gcc，go的就更快呢

最后找出原因，cgo是自动开o2优化的...