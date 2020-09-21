fn foo(x:Box<dyn Fn(u8)->u8>)->Vec<u8>{
    vec![1,2,3,4].into_iter().map(x).collect()
}

fn main(){
    println!("{:?}",foo(Box::new(|x| x+1)))
}
