struct One(i32);

impl One {
    fn call(&self,i:i32,f:impl Fn(i32)){
        f(i)
    }

    fn apply(&self,i:i32){
        self.call(i,|i|self.add(i))
    }

    fn add(&self,i:i32){
        println!("{}",i+self.0)

    }
}

fn main(){
    let one =One(1);
    one.apply(5)
}
