#include <iostream>

using namespace std;
//基类


/* 
另外多继承(环状继承),A->D, B->D, C->(A，B)，例如：

class D{......};
class B: public D{......};
class A: public D{......};
class C: public B, public A{.....};
这个继承会使D创建两个对象,要解决上面问题就要用虚拟继承格式

格式：class 类名: virtual 继承方式 父类名

class D{......};
class B: virtual public D{......};
class A: virtual public D{......};
class C: public B, public A{.....};
虚继承--（在创建对象的时候会创建一个虚表）在创建父类对象的时候

A:virtual public D
B:virtual public D
 */
class D{
    public:
        D(){cout<<"D()"<<endl;}
        ~D(){cout<<"~D()"<<endl;}
    protected:
        int d;
};

class B:virtual public D{
    public:
        B(){cout<<"B()"<<endl;}
        ~B(){cout<<"~B()"<<endl;}
    protected:
        int b;
};

class A:virtual public D{
    public:
        A(){cout<<"A()"<<endl;}
        ~A(){cout<<"~A()"<<endl;}
    protected:
        int a;
};

class C:public B, public A{
    public:
        C(){cout<<"C()"<<endl;}
        ~C(){cout<<"~C()"<<endl;}
    protected:
        int c;
};

int main(){
    cout << "Hello World!" << endl;
    C c;   //D, B, A ,C
    cout<<sizeof(c)<<endl;
    return 0;
}