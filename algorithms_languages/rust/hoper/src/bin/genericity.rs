#![feature(type_ascription)]
fn main() {
    let a = A {};
    //let a = struct{}{};  不能使用匿名struct
    //let a = struct B{}{};一样不行
    let b = ();
    let c = println!("a:{:?},b:{:?},c:{:?}",a,b,get_result_true(5).unwrap());
    println!("c:{:?}",c);
    let mut d = B::new();
    d.edit(6);
    println!("d:{:?}",d);
}

#[derive(Debug)]
struct A {}

#[derive(Debug)]
struct B<T>{
    //T 不能匿名字段
    t:T
}

//原来是tm压根编译不过去 mismatched types
/*fn get_result<T>()->Result<B<T>,()>{
    Ok(B{t:5})
}

fn get_result_err<T>()->Result<T,()>{
    Ok(5)
}*/

fn get_result_true<T>(t:T)->Result<T,()>{
    Ok(t)
}

impl B<i32>{
    fn new() -> Self{
        B{t:5}
    }
}

impl<T> B<T>{
    fn edit(&mut self,t:T){
        self.t= t;
    }
}
