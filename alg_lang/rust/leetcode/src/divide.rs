///两数相除

pub fn divide(dividend: i32, divisor: i32) -> i32 {
    if dividend == -0x80000000 && divisor == -1 {
        return 0x7fffffff;
    }
    if dividend == 0 {
        return 0;
    }
    if divisor == 1 { return dividend; }
    if divisor == -1 { return -dividend; }
    let mut result = 0;
    let mut dividend_copy = dividend;
    let mut divisor_copy = divisor;
    if dividend_copy > 0 { dividend_copy = -dividend_copy }
    if divisor_copy > 0 { divisor_copy = -divisor_copy }
    while dividend_copy <= divisor_copy {
        dividend_copy = dividend_copy - divisor_copy;
        result += 1;
    }
    if (dividend > 0 && divisor < 0) || (dividend < 0 && divisor > 0) {
        return -result;
    }
    result
}