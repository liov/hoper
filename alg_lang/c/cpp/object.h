#include <iostream>
using namespace std;

class Object {
public:
    Object()= default;
    explicit Object( double len );             // 简单的构造函数
    Object( const Object &obj);      // 拷贝构造函数
    explicit Object(double len,double het);
    ~Object(){
        cout << "~Box" << endl;
    }
private:
    double length{};      // 长度
    double breadth{};     // 宽度
    double height{};      // 高度
    double *ptr{};
};

Object::Object(double len) {
    length = len;
    ptr = new double;
    *ptr = len;
}

Object::Object(const Object &obj) {
    cout << "调用拷贝构造函数并为指针 ptr 分配内存" << endl;
    ptr = new double;
    *ptr = *obj.ptr; // 拷贝值
}

Object::Object(double len, double het) {
    length =len;
    ptr = new double;
    *ptr = len;
    height = het;
}
