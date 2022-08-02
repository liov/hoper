-fvisibility=default|internal|hidden|protected

gcc的visibility是说，如果编译的时候用了这个属性，那么动态库的符号都是hidden的，除非强制声明。

1.创建一个c源文件，内容简单
```c
#include<stdio.h>
#include<stdlib.h>


__attribute ((visibility("default"))) void not_hidden ()
{
printf("exported symbol/n");
}

void is_hidden ()
{
printf("hidden one/n");
}
```
想要做的是，第一个函数符号可以被导出，第二个被隐藏。
先编译成一个动态库，使用到属性-fvisibility
gcc -shared -o libvis.so -fvisibility=hidden vis.c

 
 动态库(.so)隐藏函数名
 一、偶遇 error: undefined reference to  xxx 问题
 
 　　尝试封装通用的接口到一个private.so，然后供客户端使用，private.so编译出来后由sample.cpp依赖调用其中封装的接口，但是一直报error: undefined reference to  xxx的错误，并且检查so、头文件都依赖正确，c方式编译的函数也用extern "C" 声明。
 
```c
 #ifdef __cplusplus
 extern "C" {
 #endif
 
  xxx
 
 #ifdef __cplusplus
 }
 #endif
```
于是用如下方法查看so的符号表根本找不到定义的 xxx 函数：
 
 readelf -s private.so
 nm -D private.so

```c
 #ifdef __cplusplus
 extern "C" {
 #endif
 
  __attribute__((visibility ("default"))) xxx
 
 #ifdef __cplusplus
 }
 #endif
```

 二、动态库函数隐藏技巧
 
 　　向客户提供动态链接库(.so)时，有些关键的函数名不希望暴露出去，此时便可以通过gcc的-fvisibility=hidden选项对编译生成的so进 的行函数符号隐藏，如：LOCAL_CPPFLAGS +=-fvisibility=hidden，执行编译后，使用nm -D xxx.so命令或者readelf --symbols xxx.so查看函数名的确被隐藏，但此时是将所有函数名都隐藏了，那么客户加载so时需要调用的接口函数名(xxx)也会找不到定义，导致编译报undefined reference to xxx错误，所以需要暴露（导出）的函数前应该增加属性__attribute__ ((visibility("default")))设置成可见。
 
 例如：
```c
 __attribute__ ((visibility("default")))  
 void hello(void)  
 {  
 }  
```
 　　实际项目开发中可以通过宏来控制，更加方便：
 
```c
 #ifdef DCIRDLL_EXPORTS
 #ifdef PLATFORM_LINUX
 #define MYDCIR_API __attribute__((visibility ("default")))  //Linux动态库(.so)
 #else
 #define MYDCIR_API __declspec(dllexport)　　　　　　　　　　　 //Windows动态库(.dll)
 #endif
 #else
 #define MTDCIR_API
 #endif
```
 　　在头文件中声明即可：
 ```c
 MTDCIR_API const voide hello();
```