let a = 5;

if (a&&1 === 0) console.log("偶数"); else console.log("奇数")
/*
能否被2^n整除 (a&&(n-1)) === 0 ? 可以:不可以
 */

console.log(a>>1)