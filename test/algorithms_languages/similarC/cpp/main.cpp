#include <iostream>

void ptr_swap(int *&v1,int *&v2){
    int *tmp =v2;
    v2=v1;
    v1=tmp;
}

int main() {
    int i =10;
    int j =20;
    int *pi=&i;
    int *pj=&j;
    std::cout<<"Before swap:\t*pi:"<<*pi<<"\t*pj:"<<*pj<<std::endl;
    ptr_swap(pi,pj);
    std::cout<<"Before swap:\t*pi:"<<*pi<<"\t*pj:"<<*pj<<std::endl;
    return 0;
}

