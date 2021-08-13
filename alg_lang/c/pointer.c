#include <stdio.h>

int main(int _argc,char* argv[]){
    int a = 10;
    int* addr_a = &a;
    printf("Address of a variable: %p\n", addr_a  );
    int* addr_b = addr_a + 1;
    printf("Address of b variable: %p\n", addr_b  );
    void* p1 = addr_a;
    void* p2 = p1 + 1;
    printf("p1: %p p2: %p\n", p1,p2  );
    int* c = (int*)p1;
    int* d = (int*)p2;
    printf("p1 + p2: %d\n", *c + *d  );
    return 0;
}