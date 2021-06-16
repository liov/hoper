#include <stdio.h>
#include <float.h>
#include <string.h>

typedef union {
   int i;
   float f;
   char  str[8];
} Foo;

typedef enum {
      MON=1, TUE, WED, THU, FRI, SAT, SUN
} DAY;

int main(int argc, char *argv[]){
    Foo foo;
    foo.i = 10;
    foo.f = 220.5f;
    strcpy( foo.str, "hello");
 
    printf( "data.i : %d\n", foo.i);
    printf( "data.f : %f\n", foo.f);
    printf( "data.str : %s\n", foo.str);

}
