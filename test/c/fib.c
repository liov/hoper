#include <time.h>

long long int fibonacci(int n);

long long int fibonacci(int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

void main(int argc, char** argv){
        int n;
        sscanf(argv[1], "%d", &n);
        fib(n);
        foo();
}

void fib(int n){
    int begintime,endtime;
    begintime = clock();
    long long int i = fibonacci(n);
    endtime = clock();
    printf("fib:%d,Running Time:%f ms\n",i, (double)(endtime-begintime));
}

void foo(){
    int begintime,endtime;
    long long int j = 0;
    begintime = clock();
   for(long long int i = 0;i<100000000;i++){
        j++;
   }
    endtime = clock();
    printf("fib:%d,Running Time:%f ms\n",j, (double)(endtime-begintime));
}
