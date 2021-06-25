#include <stdint.h>

typedef void (*rust_callback)(int32_t);

void run_callback(int32_t data, rust_callback callback) {
    callback(data); // 调用传过来的回调函数
}

unsigned long long int fibonacci(unsigned int n) {
 if (n < 2) return 1;
    return fibonacci(n - 1) + fibonacci(n - 2);
}
