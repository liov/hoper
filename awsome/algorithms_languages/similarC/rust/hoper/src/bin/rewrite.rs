//这种写法可以命名相同函数，测试和实际调用的是不同函数

#[cfg(not(test))]
fn clone() -> i32 {
    5
}

#[cfg(test)]
fn clone() -> i32 {
    10
}

fn main(){
    println!("{}",clone())
}
