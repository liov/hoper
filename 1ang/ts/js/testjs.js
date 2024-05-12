var items = [ 1, 2, 3, 4, 5, 6 ];
var results = [];
var running = 0;
var limit = 2;

function async(arg, callback) {
    console.log('参数为 ' + arg +' , 1秒后返回结果');
    setTimeout(function () { callback(arg * 2); }, 1000);
}

function final(value) {
    console.log('完成: ', value);
}

function launcher() {
    while(running < limit && items.length > 0) {
        var item = items.shift();
        async(item, function(result) {
            results.push(result);
            running--;
            if(items.length > 0) {
                launcher();
            } else if(running == 0) {
                final(results);
            }
        });
        running++;
    }
}

function bubbleSort(...values){
    for (let i=0; i<values.length-1; i++){
        for (let j=0; j<values.length-1-i; j++){
            if (values[j]>values[j+1]){
                let temp =values[j];
                values[j] = values[j+1];
                values[j+1] =temp;
            }
        }
    }
    return values;
}

let first = [11, 22];
let second = [3, 4];
let bothPlus = [20, ...first, ...second, 15];
second = [5,6];
console.log(bothPlus,second)

console.log(bubbleSort(...bothPlus))