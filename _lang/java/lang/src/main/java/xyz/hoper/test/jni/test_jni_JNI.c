#include "test_jni_JNI.h"
#include <time.h>

//gcc test_jni_JNI.c -shared -o /f/tool/dll/hello.dll -I /e/jdk-12/include -I /e/jdk-12/include/win32 -O2

long long int fibonacci(int n);

long long int fibonacci(int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

JNIEXPORT void JNICALL Java_test_jni_JNI_testHelloVoid (JNIEnv *env, jobject obj) {
  puts("hello world return void");
}

JNIEXPORT jstring JNICALL Java_test_jni_JNI_testHello (JNIEnv *env, jobject obj) {
    const char *p ="hello world return jstring";
  return (*env)->NewStringUTF(env,p);
}

JNIEXPORT jlong JNICALL Java_test_jni_JNI_fib (JNIEnv *env, jclass clazz, jint n){
        int begintime,endtime;
        begintime = clock();
        long long int i = fibonacci(n);
        endtime = clock();
        printf("Running Time:%f ms\n", (double)(endtime-begintime));
        return i;
}


JNIEXPORT void JNICALL Java_test_jni_JNI_jnifor (JNIEnv *env, jclass clazz){
       long long int j;
       for(long long int i = 0;i<100000000;i++){
            j++;
       }
}
