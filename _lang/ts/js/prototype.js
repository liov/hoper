let a = {};
let b = {x: 1};
Object.setPrototypeOf(a, b);

a.x=2
console.log(a,b)



/*let F = function () {
    this.foo = 'bar';
};

let f = new F();
// 等同于
let f = Object.setPrototypeOf({}, F.prototype);

F.call(f);*/



// 原型对象
let A = {
    a:1,
    print: function () {
        console.log(this.a);
    }
};

// 实例对象
let B = Object.create(A);

A.a=2;
B.print() // hello
console.log(A,B)
console.log(Object.getPrototypeOf(B))



let obj = Object.create({}, {
    p1: {
        value: 123,
        enumerable: true,
        configurable: true,
        writable: true,
    },
    p2: {
        value: 'abc',
        enumerable: true,
        configurable: true,
        writable: true,
    }
});

console.log(obj)



function copyObject(orig) {
    return Object.create(
        Object.getPrototypeOf(orig),
        Object.getOwnPropertyDescriptors(orig)
    );
}

let a1 = {x:1,b:2};
let b1 = copyObject(a1);
console.log(a1,b1);
