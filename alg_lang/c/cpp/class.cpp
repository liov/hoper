#include <iostream>
#include "box.h"
using namespace std;

int main(void){
   Rectangle Rect;
   Triangle  Tri;
   Shape *shape;

   Rect.setWidth(5);
   Rect.setHeight(7);
   // 输出对象的面积
   cout << "Total Rectangle area: " << Rect.getArea() << endl;
   Rectangle Rect2 = Rect(3,5,7);
   shape = &Rect2;
    cout << "area: " << shape->getArea() << endl;

   Tri.setWidth(5);
   Tri.setHeight(7);
   // 输出对象的面积
   cout << "Total Triangle area: " << Tri.getArea() << endl; 
 
   return 0;
}