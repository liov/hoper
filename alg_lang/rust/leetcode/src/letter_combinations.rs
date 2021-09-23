///电话号码的字母组合

//Y 组合子是 lambda 演算中的一个概念，是任意函数的不动点，在函数式编程中主要作用是 提供一种匿名函数的递归方式。
//"λ" 字符可以看作 function 声明，"."字符前为参数列表，"."字符后为函数体。
//不动点是函数的一个特征：对于函数 f(x)，如果有变量  a 使得  f(a)=a 成立，则称 a 是函数 f 上的一个不动点。

//注意结尾没有分号的那一行 x+1，与你见过的大部分代码行不同。表达式的结尾没有分号。如果在表达式的结尾加上分号，
// 它就变成了语句，而语句不会返回值。在接下来探索具有返回值的函数和表达式时要谨记这一点。

pub fn letter_combinations(digits: String) -> Vec<String> {
    let combination = vec![
        vec![b'a', b'b', b'c'],
        vec![b'd', b'e', b'f'],
        vec![b'g', b'h', b'i'],
        vec![b'j', b'k', b'l'],
        vec![b'm', b'n', b'o'],
        vec![b'p', b'q', b'r', b's'],
        vec![b't', b'u', b'v'],
        vec![b'w', b'x', b'y', b'z'],
    ];

    let mut result = Vec::new();

    let cap = digits.len();
    let mut middle: Vec<u8> = Vec::with_capacity(cap);
    unsafe { middle.set_len(cap); }

    let digits_vec = digits.as_bytes();
    //闭包递归实在不会，改成了函数递归，消耗几何未知，有个clone()难受啊
    fn get(n: usize, res: &mut Vec<String>, mid: &mut Vec<u8>, com: &Vec<Vec<u8>>, dig: &[u8]) {
        if n < mid.len() {
            for i in &(com[dig[n] as usize - 50]) {
                mid[n] = *i;
                get(n + 1, res, mid, com, dig);
                if n == mid.len() - 1 {
                    res.push(String::from_utf8(mid.clone()).unwrap());
                }
            }
        }
    }


    //let Y = |y|(|x|y(x(x)))(|x|y(|n|x(x)(n)));

    get(0, &mut result, &mut middle, &combination, &digits_vec);
    result
}
