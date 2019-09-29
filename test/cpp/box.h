#include <iostream>
using namespace std;
/* 继承类型
当一个类派生自基类，该基类可以被继承为 public、protected 或 private 几种类型。继承类型是通过上面讲解的访问修饰符 access-specifier 来指定的。

我们几乎不使用 protected 或 private 继承，通常使用 public 继承。当使用不同类型的继承时，遵循以下几个规则：

公有继承（public）：当一个类派生自公有基类时，基类的公有成员也是派生类的公有成员，基类的保护成员也是派生类的保护成员，基类的私有成员不能直接被派生类访问，但是可以通过调用基类的公有和保护成员来访问。
保护继承（protected）： 当一个类派生自保护基类时，基类的公有和保护成员将成为派生类的保护成员。
私有继承（private）：当一个类派生自私有基类时，基类的公有和保护成员将成为派生类的私有成员。 */
class Box {
    public:
        // 纯虚函数
        virtual double getVolume() = 0;
        ~Box(){
            cout << "~Box" << endl;
        }
    private:
        double length;      // 长度
        double breadth;     // 宽度
        double height;      // 高度
};

// 基类
class Shape {
    public:
        // 提供接口框架的纯虚函数
        virtual int getArea() = 0;
        void setWidth(int w){
            width = w;
        }
        void setHeight(int h){
            height = h;
        }
        protected:
        int width;
        int height;
        int length;
};

// 基类 PaintCost
class PaintCost {
   public:
      int getCost(int area){
         return area * 70;
      }
};

/* 
可重载运算符/不可重载运算符
下面是可重载的运算符列表：

双目算术运算符	+ (加)，-(减)，*(乘)，/(除)，% (取模)
关系运算符	==(等于)，!= (不等于)，< (小于)，> (大于>，<=(小于等于)，>=(大于等于)
逻辑运算符	||(逻辑或)，&&(逻辑与)，!(逻辑非)
单目运算符	+ (正)，-(负)，*(指针)，&(取地址)
自增自减运算符	++(自增)，--(自减)
位运算符	| (按位或)，& (按位与)，~(按位取反)，^(按位异或),，<< (左移)，>>(右移)
赋值运算符	=, +=, -=, *=, /= , % = , &=, |=, ^=, <<=, >>=
空间申请与释放	new, delete, new[ ] , delete[]
其他运算符	()(函数调用)，->(成员访问)，,(逗号)，[](下标)
下面是不可重载的运算符列表：

.：成员访问运算符
.*, ->*：成员指针访问运算符
::：域运算符
sizeof：长度运算符
?:：条件运算符
#： 预处理符号 */


// 派生类
class Rectangle: public Shape, public PaintCost{
    public:
        ~Rectangle(){
            cout << "~Rectangle" << endl;
        }
        
        int getArea(){ 
            return (width * height); 
        }
        double getVolume(void){
            return length * width * height;
        }
        
        void setLength(int len ){
            length = len;
        }

        void setLength(double len ){
            length = len;
        }
    
        void setWidth( double wid ){
            width = wid;
        }
    
        void setHeight( double hei ){
            height = hei;
        }

        void toString(){
            cout << "length: " << length << "width: " << width << "height: " << height <<endl;
        }

        // 重载 + 运算符，用于把两个 Box 对象相加
        Rectangle operator+(const Rectangle& r){
            Rectangle rectangle;
            rectangle.length = this->length + r.length;
            rectangle.width = this->width + r.width;
            rectangle.height = this->height + r.height;
            return rectangle;
        }

        Rectangle operator()(int a, int b, int c){
            Rectangle rectangle;
            rectangle.length =  a + c;
            rectangle.width = a + b;
            rectangle.height = b + c;
            return rectangle;
        }
};

class Triangle: public Shape{
    public:
        int getArea(){ 
            return (width * height)/2; 
        }
};