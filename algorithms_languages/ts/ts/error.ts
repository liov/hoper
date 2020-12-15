class Animal {
    name: string;
}
class Dog extends Animal {
    breed: string;
}

// 错误：使用数值型的字符串索引，有时会得到完全不同的Animal!
interface NotOkay {
    // @ts-ignore
    [x: number]: Animal;
    [x: string]: Dog;
}

let animal1:Animal = {name:"dog1"};
// @ts-ignore
let dog1:Dog ={breed:"labuladuo1"};
let animal2:Animal = {name:"dog2"};
// @ts-ignore
let dog2:Dog ={breed:"labuladuo2"};
let a:NotOkay={};
a[1]=animal1;
a[2]=animal2;
a["1"]=dog1;
a["2"]=dog2;
console.log(a[1]);
