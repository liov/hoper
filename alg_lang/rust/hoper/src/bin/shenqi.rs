//可以运行修改
fn one() {
    unsafe {
        let stk_arr = [b'a', b' ', b'v', b'a', b'r', b' ', b'o', b'f', b' ', b'a'];
        let a = std::str::from_utf8_unchecked(&stk_arr);
        let addr_a = a.as_ptr() as *mut u8;
        let c = b"c var of c";

        println!("a: {}", a);
        println!("a的地址：{:p}", addr_a);

        addr_a.copy_from(c.as_ptr(), 10);
        let addr_a = a.as_ptr() as *mut u8;
        println!("a的地址：{:p}", addr_a);
        println!("a: {}", a);
    }
}
//不可以运行修改
fn two() {
    unsafe {
        let a = std::str::from_utf8_unchecked(b"a var of a");
        let addr_a = a.as_ptr() as *mut u8;
        let c = b"c var of c";

        println!("a: {}", a);
        println!("a的地址：{:p}", addr_a);

        addr_a.copy_from(c.as_ptr(), 10);
        let addr_a = a.as_ptr() as *mut u8;
        println!("a的地址：{:p}", addr_a);
        println!("a: {}", a);
    }
}

fn main() {}