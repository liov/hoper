
//不用mod判断整数能否被2整除，位运算
function divisibility(x){
    return (x&1)===0
}

console.log(divisibility(3))
console.log(divisibility(4))
