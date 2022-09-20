///字符串相乘

//执行用时 : 0 ms, 在Multiply Strings的Rust提交中击败了100.00% 的用户
//内存消耗 : 2 MB, 在Multiply Strings的Rust提交中击败了100.00% 的用户
/*m 位数和n位数相乘，结果位数为m+n-1或m + n，因此存进位数据数组大小申请为m+n位
对应位相乘后的结果与进位数据数组对应位置相加，十位数存入进位数组下一位，个位数留在该位
存进位数据数组在转为字符串返回的时候，要把前导零去掉*/
pub fn multiply(num1: String, num2: String) -> String {
    let num_vec1 = num1.as_bytes();
    let num_vec2 = num2.as_bytes();
    let cap = num_vec1.len() + num_vec2.len();
    let mut result = Vec::with_capacity(cap);
    unsafe { result.set_len(cap); }
    let con = 48;
    let mut product; //乘积
    let mut decade = 0; //十位

    if (num_vec1.len() == 1 && num_vec1[0] == con) || (num_vec2.len() == 1 && num_vec2[0] == con) {
        return String::from("0");
    }

    for i in 0..num_vec1.len() {
        for j in 0..num_vec2.len() {
            if i == 0 { result[cap - i - j - 1] = 0 } else if j == 0 { result[cap - i - num_vec2.len() - 1] = 0 }
            product = (num_vec1[num_vec1.len() - i - 1] - con) * (num_vec2[num_vec2.len() - j - 1] - con) + decade;
            decade = product / 10 + (product % 10 + result[cap - i - j - 1]) / 10;
            result[cap - i - j - 1] = (product % 10 + result[cap - i - j - 1]) % 10;

            if i == num_vec1.len() - 1 {
                result[cap - i - j - 1] = result[cap - i - j - 1] + con
            } else if j == 0 {
                result[cap - i - j - 1] = result[cap - i - j - 1] + con
            }
        }
        result[cap - i - num_vec2.len() - 1] = decade;
        decade = 0;
    }

    if result[0] != 0 { result[0] = result[0] + con } else { result.remove(0); }
    String::from_utf8(result).unwrap()
}