function f(x,y){//定义一个简单的函数
    return x+y;
}
function o(a,b){//定义一个函数结构体
    return a*b;
}
console.log(f.call(o,3,4))//返回值为7
console.log(f.apply(o,[3,4]))//返回值为7
//call和apply方法是把一个函数转化为指定对象的方法，并在该对象上调用该函数，函数并没有做为对象的方法而存在，函数被动态调用之后，临时方法就会被注销。
function f1(){console.log("调用f1",this.a)}//定义空函数
const object ={a:1}
f1.call(object);//把函数f()绑定为object方法
//call和apply方法可以动态的改变函数this指代的对象。下面使用call方法不断改变函数内this指代对象，通过改变call方法的第一个参数来实现
const x="o"
function a(){
    this.x="a";
}
function c(){
    console.log(x);
}
function d(){
    console.log(this.x);
}

c.call(new a());//返回字符a，及函数a内部局部变量x的值，this此时指向函数a
d.call(new a());//返回字符b，及函数b内部局部变量心得值，this此时指向函数b
d.call(c);//返回undefinded，及函数d内部的函数变量x的值，但该函数并没有定义x变量所以返回无定义，this此时指向函数c
