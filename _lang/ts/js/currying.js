// 普通的add函数
function add(x, y) {
    return x + y
}

// Currying后
function curryingAdd(x) {
    return function (y) {
        return add(x, y)
    }
}

add(1, 2)           // 3
curryingAdd(1)(2)   // 3

//柯里化，英语：Currying(果然是满满的英译中的既视感)，
// 是把接受多个参数的函数变换成接受一个单一参数（最初函数的第一个参数）的函数，
// 并且返回接受余下的参数而且返回结果的新函数的技术。
