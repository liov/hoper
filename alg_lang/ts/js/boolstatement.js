const readline = require('readline');

console.log(process.argv);

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

/*let inputArr = [];
rl.on('line', function (input) {
    inputArr = input.split(" ");
    inputArr.forEach(function(item,index){
        inputArr[index] = +item;// 转化为数字
    });
    console.log(inputArr);
    inputArr = [];// 清空数组
    rl.close();
});*/

rl.on('line', function (input) {
    input==='close'&&rl.close(); //if语句,&&阻断
    let num = +input;
    console.log(num);
});

rl.on('close', function() {
    console.log('程序结束');
    process.exit(0);
});

