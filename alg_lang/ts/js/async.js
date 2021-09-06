//异步函数加了await相当于是同步的，但是如果其他函数掉这个异步函数没有await是异步执行的
const fs = require('fs');

function writeFile(fileName,data) {
    return new Promise(function (resolve, reject) {
      fs.writeFile(fileName,data, function(error) {
        if (error) return reject(error);
        resolve(fileName+'写入成功');
      });
    });
  };

async function foo(){
   let r = await writeFile('../foo.txt','foo')
   console.log(r)
}

async function writeFileAsync(){
    await foo()
    console.log('正在异步写入')
}

writeFileAsync()

const readFile = function (fileName) {
    return new Promise(function (resolve, reject) {
      fs.readFile(fileName, function(error, data) {
        if (error) return reject(error);
        resolve(data);
      });
    });
  };

async function bar(){
    let f = await readFile('../foo.txt')
    console.log(f.toString())
}

function writeFileSync(){
    bar()
    console.log('正在同步写入')
}

writeFileSync()

function fib(n){
  return n<2?1:fib(n-2)+fib(n-1)
}

async function fibAsync(n){
  return new Promise(function (resolve, reject){
    resolve(fib(n)); 
  })
}

async function singleThread(n){
  console.log('计算斐波那契')
  let data = await fibAsync(n)
  console.log(data)
}
//这个函数证明了什么，如我所想，node是单线程的，所以在涉及非io操作时，异步函数是在一个线程里执行的
//本应在计算的同时就输出计算中，但是是计算完成后才输出的
function singleThreadTest(){
  singleThread(30)
  console.log('计算中')
}

singleThreadTest()

function longTime(path){
  return new Promise(function (resolve, reject) {
    const file_size = 1024*1024*1024;         //1G  
    const buf_size = 10240;
    let buf = Buffer.alloc(buf_size);
    let temp = Buffer.alloc(10);
    for (let i=0; i<10; i++) {
        temp[i] = i.toString().charCodeAt(0);  
    }
  // init buf
  for (var i=0; i<buf_size/10-1; i++) {  
    temp.copy(buf, 10*i);
  }
  temp.copy(buf, 10*i, 0, buf_size-parseInt(buf_size/10)*10);
  // write to file
  fs.open(path, 'w', 0o666, function(err, fd){
    if (err) throw reject(err);
    var i=0;
    function write(err, written) {
        if (err) throw reject(err);
        if (i>=file_size/buf_size) {    //close the file  
            fs.close(fd,err=> {if (err) throw reject(err)});
            resolve('完成');
        } else {            //continue to write
            let length = buf_size;  
            if ((i+1)*buf_size>file_size) {
                length = file_size-i*buf_size;
            }
            fs.write(fd, buf, 0, length, null, write);
            i++;
        }
    }
    write(null, 0);
  });
});
}
//每个await都会等待知道异步操作完成继续向下执行
async function twoAwati(){
  let data = await longTime('../foo1.txt')
  let r = await writeFile('../foo2.txt','foo')
  console.log(r)
  console.log(data)
}

twoAwati()