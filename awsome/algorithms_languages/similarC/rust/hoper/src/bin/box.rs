fn main() {
    let a = 5;
    let add_a = &a as *const i32;
    println!("{}", unsafe { *add_a });
    let b = Box::new(5);
    let addr_b = &*b as *const i32;
    println!("{:p}:{:p}", *&b, addr_b);
    test1(b); //注释此行，获取正确的值，不注释，错误的值
    println!("{:p}:地址addr_b的值：{}", addr_b, unsafe { *addr_b });

    let c = Some(User { name: String::from("jyb"), age: 16 });
    let addr_c = c.as_ref().unwrap();
    println!("{:p},{:?}", addr_c, *addr_c);
    let addr_cc = &c.unwrap();
    println!("{:p},{:?}", addr_cc, *addr_cc);
    //test3(c);

    let d = Box::new(5);
    print!("{:p}:{:p}", &d, d.as_ref())
}

fn test1(a: Box<i32>) {
    let addr = a.as_ref();
    println!("{:p}:{}", addr, *addr);
    //此处获取所有权a，并执行完，释放了内存，因此最后获取到错误的值
}

fn test3(a: Option<User>) {
    let addr = a.as_ref().unwrap();
    println!("{:p}:{:?}", addr, *addr);
    //此处获取所有权a，并执行完，释放了内存，因此最后获取到错误的值
}

fn test2(a: i32) {
    println!("{:p}", &a);
}

#[derive(Debug)]
struct User {
    name: String,
    age: i32,
}
