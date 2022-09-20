function PointES5(x,y){
    this.x=x;
    this.y=y;
}
PointES5.prototype.toString=function(){
    return '('+this.x+', '+this.y+')';
}
var point=new PointES5(1,2);


class Point {
    constructor(x, y) {
        this.x = x;
        this.y = y;
    }

    toString() {
        return '(' + this.x + ', ' + this.y + ')';
    }
}

// 工厂模式
function Car1(name, color, price) {
    var tempcar = new Object;
    tempcar.name = name;
    tempcar.color = color;
    tempcar.price = price;
    tempcar.getCarInfo = function () {
        console.log(`name: ${this.name},color: ${this.color},price: ${this.price}`);
    }
    return tempcar;
}

var mycar1 = new Car1('BMW', 'red', '100000');
mycar1.getCarInfo();

//缺点：每次 new 一个对象的时候，都会重新创建一个 getCaeInfo() 函数；

// 构造函数方式
function Car2(name, color, price) {
    this.name = name;
    this.color = color;
    this.price = price;
    this.getCarInfo = function () {
        console.log(`name: ${this.name},color: ${this.color},price: ${this.price}`);
    }
}

var myCar2 = new Car2('桑塔纳', 'green', '123456');
myCar2.getCarInfo();

/*优点：
不用创建临时对象；
不用返回临时对象；
缺点：与‘工厂模式’相同,重复创建函数；*/

// 原型方式
function Car3(name, color, price) {
    Car.prototype.name = name;
    Car.prototype.color = color;
    Car.prototype.price = price;
    Car.prototype.getCarInfo = function () {
        console.log(`name: ${this.name},color: ${this.color},price: ${this.price}`);
    }
}

var myCar3 = new Car3('兰博基尼', 'red', '10000000000');
myCar3.getCarInfo();

/*优点：
解决了重复创建函数的问题；
可以使用 instanceof 检验类型 myCar instanceof Car // true
缺点：
多个实例创建的相同属性指向同一块内存；
例子：*/
Car.prototype.drivers = ['Tim', 'Jone'];
myCar3.drivers.push('mike');
console.log(myCar3.drivers); // ['Tim', 'Jone', 'mike']

// 动态原型方式（推荐）
function Car4(name, color, price, drivers) {
    this.name = name;
    this.color = color;
    this.price = price;
    this.drivers = drivers;
}

Car4.prototype.getCarInfo = function () {
    console.log(`name: ${this.name},color: ${this.color},price: ${this.price}`);
}

var myCar4 = new Car4('兰博基尼', 'red', '10000000000', ['qaz', 'wsx']);
myCar4.drivers.push('mi');
console.log(myCar4.drivers); // ["qaz", "wsx", "mi"]

var myCar5 = new Car4('兰博基尼1', 'red1', '100000000001', ['qaz1', 'wsx1']);
myCar5.drivers.push('mi1');
console.log(myCar5.drivers); // ["qaz1", "wsx1", "mi1"]

/*思想：
类的属性 要随实例对象动态改变； => 动态
类的方法 要随原型保持不变；=> 原型*/
