#![feature(async_await, futures_api)]

use hoper::utils::tree::MyTree;
use std::ops::Deref;


fn main() {

    let x = std :: f64 :: consts :: PI;
    let mut t = MyTree::new();
    t.insert(x);
    println!("{:?}",t);
    t.insert(3f64);
    println!("{:?}",t);
    t.insert(5f64);
    println!("{:?}",t);
    t.insert(9f64);
    println!("{:?}",t);
    t.insert(66f64);
    println!("{:?}",t);
    t.insert(18f64);
    println!("{:?}",t);
    t.insert(12f64);
    println!("{:?}",t);
    t.insert(111f64);
    println!("{:?}",t);
    t.insert(12f64);
    println!("{:?}",t);
    t.insert(2f64);
    println!("{:?}",t);
    t.peek();

}



struct MyBox<T>(T);
impl<T> MyBox<T> {
    fn new(x: T) -> MyBox<T> {
        MyBox(x)
    }
}
impl<T> Deref for MyBox<T> {
    type Target = T;

    fn deref(&self) -> &T {
        &self.0
    }
}
impl<T> Drop for MyBox<T> {
    fn drop(&mut self) {
        println!("清理Mybox")
    }
}
fn hello(name: &str) {
    println!("Hello, {}!", name);
}

#[derive(Debug)]
enum Message {
   Quit,
   Move { x: i32, y: i32 },
   Write(String),
   ChangeColor(i32, i32, i32),
}

impl Message {
    async  fn call(&self) {
       // 在这里定义方法体
       async {
            // 省略业务代码
            "set state".to_owned()
        };
   }
}

trait Conn {
   fn connect(&self) ->i32;
}

struct Cacher<T>
   where T: Fn(u32) -> u32
{
   calculation: T,
   value: Option<u32>,
}

impl<T> Cacher<T> where T: Fn(u32) -> u32 {
   fn new(calculation: T) -> Cacher<T> {
       Cacher {
           calculation,
           value: None,
       }
   }

   fn value(&mut self, arg: u32) -> u32 {
       match self.value {
           Some(v) => v,
           None => {
               let v = (self.calculation)(arg);
               self.value = Some(v);
               v
           },
       }
   }
}
