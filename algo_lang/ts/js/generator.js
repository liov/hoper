function* gen(x){
    //函数在此终止执行,返回yield后的值
    yield x + 1;
    console.log("执行了")
    var y = yield x + 2;
    return y;
}

var g = gen(1);
console.log("执行前后")
console.log(g.next()) // { value: 3, done: false }
console.log(g.next(2)) // { value: 2, done: true }
console.log(g.next(3)) // { value: 2, done: true }
//生成器虽然是异步的，但是跟async函数还是有区别的

//迭代器
let obj = {
    * [Symbol.iterator]() {
      yield 'hello';
      yield 'world';
    }
  };
  
  for (let x of obj) {
    console.log(x);
  }

//模拟生成器，内部保存一个状态
function idMaker() {
    var index = 0;

    return {
        next: function() {
        return {value: index++, done: false};
        }
    };
}

var it = idMaker();

it.next().value // 0
it.next().value // 1