//期望引用自身
#[derive(Debug)]
struct Foo<'a,T> {
    val:Vec<T>,
    val_ref:Vec<&'a T>
}

impl<'a,T> Foo<'a,T> {
    fn new() -> Self {
        Foo{
            val: Vec::new(),
            val_ref: Vec::new()
        }
    }
    //这个方法是不可行的，val_ref永久的持有了val的可变引用
    //理论上来讲，结构体是成立的，都放的是不可变引入，但操作不能这么实现
    //RC or unsafe
    //真的是动不动就要unsafe，cpp大法好
    fn push(&'a mut self,new_val:T){
        self.val.push(new_val);
        self.val_ref.push(&(self.val[self.val.len()-1]));
    }

}

struct Bar<T> {
    val:T,
    val_ref:*mut T,
}

fn main(){
    let mut a = Foo::<i32>::new();
    //a.push(1);
    println!("{:?}",a);
}
