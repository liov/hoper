//rust 枚举是函数？
#[derive(Debug)]
enum A {
    A1(i32, i32),
    A2(String),
    A3(u8),
    A4(u128),
}

type B = (i32,i32);
#[derive(Debug)]
struct C(i32,i32);

fn foo(a: i32, b: i32) -> A {
    A::A1(a, b)
}

fn bar(a: i32, b: i32)  {
    a + b;
}

fn main() {
    let a = A::A1(0, 0);
    let b: fn(i32, i32) -> A = foo;
    let c: fn(i32, i32) -> A = A::A1;
    let d = c;
    println!("{:?}", d(1, 2));
    //let e:fn(i32,i32) = B;//expected value, found type alias `B`
    let f:fn(i32,i32)->C = C;
    println!("{:?}", f(1, 2));
}
